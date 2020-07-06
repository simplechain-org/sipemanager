package controllers

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"
	"sipemanager/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type ErrLogCode struct {
	message string
	code    int
	err     string
}

func GetRpcApi(node dao.InstanceNodes) (*blockchain.Api, error) {
	n := &blockchain.Node{
		Address:   node.Address,
		Port:      node.Port,
		ChainId:   node.ChainId,
		IsHttps:   node.IsHttps,
		NetworkId: node.NetworkId,
	}
	api, err := blockchain.NewApi(n)
	if err != nil {
		return nil, err
	}
	return api, nil

}

func (this *Controller) ListenCrossEvent() {
	cron := cron.New()
	cron.AddFunc("@every 5s", func() {
		fmt.Println("current event time is ", time.Now())
		nodes, err := this.dao.GetInstancesJoinNode()
		filterNodes := utils.RemoveRepByLoop(nodes)
		if err != nil {
			logrus.Error(&ErrLogCode{message: "routers => ListenEvent:", code: 30002, err: "cant not found nodes"})
		}
		fmt.Printf("-------nodes-----%+v\n", filterNodes)
		go this.createCrossEvent(nodes)
	})
	cron.Start()
}

func (this *Controller) ListenBlock() {
	var group sync.WaitGroup
	NodeChannel := make(chan BlockChannel)
	fmt.Println("current event time is ", time.Now())
	nodes, err := this.dao.GetInstancesJoinNode()
	filterNodes := utils.RemoveRepByLoop(nodes)
	//count := len(filterNodes)
	if err != nil {
		logrus.Warn(&ErrLogCode{message: "routers => ListenEvent:", code: 30001, err: err.Error()})
	}
	go this.createBlock(filterNodes, &group, NodeChannel)

	for i := 0; i <= len(filterNodes); i++ {
		ch, ok := <-NodeChannel
		logrus.Infof("node channel is %+v, ok = %+v", ch, ok)
		if ok {
			go this.HeartChannel(ch, group, NodeChannel)
		}
	}
	//for range NodeChannel {
	//	count--
	//	if count == 0 {
	//		close(NodeChannel)
	//	}
	//}

}

func (this *Controller) createCrossEvent(nodes []dao.InstanceNodes) {
	a := time.Now()
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("current nodes %+v ", nodes[i])
		contract, err := this.dao.GetContractById(nodes[i].ContractId)
		blockNumber := this.dao.GetMaxCrossNumber(nodes[i].ChainId)
		addresses := []common.Address{
			common.HexToAddress(nodes[i].CrossAddress),
		}

		records := ethereum.FilterQuery{
			FromBlock: big.NewInt(blockNumber),
			Addresses: addresses,
		}
		api, err := GetRpcApi(nodes[i])
		logs, err := api.GetPastEvents(records)
		if err != nil {
			logrus.Errorf("FilterLogs:%v", err)
		}

		abiParsed, err := abi.JSON(strings.NewReader(contract.Abi))
		if err != nil {
			logrus.Warn(err.Error())
		}

		this.EventLog(logs, abiParsed, nodes[i])

	}
	fmt.Println(time.Since(a))
}

type CrossMakerTx struct {
	TxId          [32]byte
	From          common.Address
	To            common.Address
	RemoteChainId *big.Int
	Value         *big.Int
	DestValue     *big.Int
	Data          []byte
	Raw           types.Log // Blockchain specific contextual infos
}

type CrossTakerTx struct {
	TxId          [32]byte
	To            common.Address
	RemoteChainId *big.Int
	From          common.Address
	Value         *big.Int
	DestValue     *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

type CrossMakerFinish struct {
	TxId [32]byte
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

func (this *Controller) EventLog(logs []types.Log, abiParsed abi.ABI, node dao.InstanceNodes) {
	makerTx := abiParsed.Events["MakerTx"].ID.Hex()
	takerTx := abiParsed.Events["TakerTx"].ID.Hex()
	makerFinish := abiParsed.Events["MakerFinish"].ID.Hex()
	for _, event := range logs {
		var item dao.CrossEvents
		switch event.Topics[0].Hex() {
		case makerTx:
			var args CrossMakerTx
			err := abiParsed.Unpack(&args, "MakerTx", event.Data)
			if err != nil {
				logrus.Error("makerTx:", err.Error())
			}
			item = dao.CrossEvents{
				BlockNumber:     event.BlockNumber,
				TxId:            event.Topics[1].Hex(),
				Event:           "MakerTx",
				From:            "0x" + event.Topics[2].Hex()[26:],
				NetworkId:       node.NetworkId,
				RemoteNetworkId: args.RemoteChainId.Int64(),
				Value:           args.Value.String(),
				DestValue:       args.DestValue.String(),
				TransactionHash: event.TxHash.Hex(),
				CrossAddress:    node.CrossAddress,
				ChainId:         node.ChainId,
			}
			this.dao.MakerEventUpsert(item)
		case takerTx:
			var args CrossTakerTx
			err := abiParsed.Unpack(&args, "TakerTx", event.Data)
			if err != nil {
				logrus.Error("takerTx:", err.Error())
			}
			item = dao.CrossEvents{
				BlockNumber:     event.BlockNumber,
				TxId:            event.Topics[1].Hex(),
				Event:           "TakerTx",
				From:            strings.ToLower(args.From.Hex()),
				To:              "0x" + event.Topics[2].Hex()[26:],
				NetworkId:       node.NetworkId,
				RemoteNetworkId: args.RemoteChainId.Int64(),
				Value:           args.Value.String(),
				DestValue:       args.DestValue.String(),
				TransactionHash: event.TxHash.Hex(),
				CrossAddress:    node.CrossAddress,
				ChainId:         node.ChainId,
			}
			this.dao.TakerEventUpsert(item)
		case makerFinish:
			item = dao.CrossEvents{
				BlockNumber:     event.BlockNumber,
				TxId:            event.Topics[1].Hex(),
				Event:           "MakerFinish",
				To:              "0x" + event.Topics[2].Hex()[26:],
				NetworkId:       node.NetworkId,
				TransactionHash: event.TxHash.Hex(),
				CrossAddress:    node.CrossAddress,
				ChainId:         node.ChainId,
			}
			this.dao.MakerFinishEventUpsert(item)
		}
	}
}

func (this *Controller) HeartChannel(ch BlockChannel, group sync.WaitGroup, NodeChannel chan BlockChannel) {

	cron := cron.New()
	cron.AddFunc("@every 5s", func() {
		current, err := this.dao.GetNodeById(ch.ChainId)
		if err != nil {
			logrus.Warn(&ErrLogCode{message: "time_task => HeartChannel:", code: 20005, err: err.Error()})
		}
		currents := []dao.InstanceNodes{
			dao.InstanceNodes{
				Address:   current.Address,
				Port:      current.Port,
				IsHttps:   current.IsHttps,
				NetworkId: current.NetworkId,
				Name:      current.Name,
				ChainId:   current.ChainId,
			},
		}
		fmt.Printf("current node is %+v", currents)
		go this.createBlock(currents, &group, NodeChannel)
		ch, ok := <-NodeChannel
		logrus.Infof("node channel is %+v, ok = %+v", ch, ok)
	})
	cron.Start()
	// Heart Recursive execution
	//if ok {
	//	go this.HeartChannel(object, ch, group, NodeChannel)
	//}
}

func (this *Controller) createBlock(nodes []dao.InstanceNodes, group *sync.WaitGroup, NodeChannel chan<- BlockChannel) {
	a := time.Now()
	for _, node := range nodes {
		group.Add(1)
		go this.syncAllNodes(node, group, NodeChannel)
	}
	fmt.Println(time.Since(a))
}

func (this *Controller) syncAllNodes(node dao.InstanceNodes, group *sync.WaitGroup, NodeChannel chan<- BlockChannel) {
	api, err := GetRpcApi(node)
	chainId := node.ChainId
	header, err := api.GetHeaderByNumber()
	dbMaxNum := this.dao.GetMaxBlockNumber(chainId)
	if err != nil {
		logrus.Warn(&ErrLogCode{message: "time_task => createBlock:", code: 20001, err: err.Error()})
	}
	newBlockNumber := header.Number.Int64()

	var j int64
	for j = dbMaxNum; j <= newBlockNumber; j++ {
		from := j
		var to int64
		if j == newBlockNumber {
			to = newBlockNumber - 12
			from = newBlockNumber
		} else {
			to = newBlockNumber
		}
		this.BlocksListen(from, to, api, node)
	}
	NodeChannel <- BlockChannel{
		ChainId:     chainId,
		BlockNumber: dbMaxNum,
	}
	group.Done()
}

func (this *Controller) BlocksListen(from int64, to int64, api *blockchain.Api, node dao.InstanceNodes) error {
	var err error
	// 重新同步最近的12个区块
	if from > to {
		for i := to; i <= from; i++ {
			fmt.Printf("Resync Block Create: %+v\n", i)
			this.SyncBlock(api, i, node)
		}
	} else {
		this.SyncBlock(api, from, node)
	}

	return err
}
