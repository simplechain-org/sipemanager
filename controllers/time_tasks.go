package controllers

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/simplechain-org/go-simplechain"
	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/sirupsen/logrus"

	"sipemanager/blockchain"
	"sipemanager/dao"
	"sipemanager/utils"
)

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
			logrus.Error(utils.ErrLogCode{LogType: "controller => time_task => ListenCrossEvent:", Code: 20006, Message: "cant not found nodes", Err: nil})
		}
		go this.createCrossEvent(filterNodes)
	})
	cron.Start()
}

func (this *Controller) ListenAnchors() {
	cron := cron.New()
	cron.AddFunc("@every 10s", func() {
		this.AnalysisAnchors()
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
		logrus.Warn(utils.ErrLogCode{LogType: "controller => time_task => ListenBlock:", Code: 20007, Message: err.Error(), Err: nil})
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
	for i := 0; i < len(nodes); i++ {
		//fmt.Printf("current nodes %+v ", nodes[i])
		contract, err := this.dao.GetContractById(nodes[i].ContractId)
		blockNumber := this.dao.GetMaxCrossNumber(nodes[i].ChainId)
		addresses := []common.Address{
			common.HexToAddress(nodes[i].CrossAddress),
		}

		records := simplechain.FilterQuery{
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
}

type CrossMakerTx struct {
	TxId          [32]byte
	From          common.Address
	To            common.Address
	RemoteChainId *big.Int
	Value         *big.Int
	DestValue     *big.Int
	Data          []byte
	Raw           types.Log
}

type CrossTakerTx struct {
	TxId          [32]byte
	To            common.Address
	RemoteChainId *big.Int
	From          common.Address
	Value         *big.Int
	DestValue     *big.Int
	Raw           types.Log
}

type CrossMakerFinish struct {
	TxId [32]byte
	To   common.Address
	Raw  types.Log
}

func (this *Controller) EventLog(logs []types.Log, abiParsed abi.ABI, node dao.InstanceNodes) {
	makerTx := abiParsed.Events["MakerTx"].ID().Hex()
	takerTx := abiParsed.Events["TakerTx"].ID().Hex()
	makerFinish := abiParsed.Events["MakerFinish"].ID().Hex()
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
			logrus.Warn(utils.ErrLogCode{LogType: "controller => time_task => HeartChannel:", Code: 20005, Message: err.Error(), Err: nil})
		}
		chain, err := this.dao.GetChain(current.ChainId)
		if err != nil {
			logrus.Error(err)
			return
		}
		currents := []dao.InstanceNodes{
			dao.InstanceNodes{
				Address:    current.Address,
				Port:       current.Port,
				IsHttps:    current.IsHttps,
				NetworkId:  chain.NetworkId,
				Name:       current.Name,
				ChainId:    current.ChainId,
				ContractId: ch.currentNode.ContractId,
			},
		}
		go this.createBlock(currents, &group, NodeChannel)
		ch, ok := <-NodeChannel
		logrus.Infof("node HeartChannel is %+v, ok = %+v", ch, ok)
		//select {
		//case <-NodeChannel:
		//	fmt.Println("消费完成……………………")
		//	return
		//case <-time.After(time.Second * 5):
		//	fmt.Println("超时………………………")
		//	return
		//}
	})
	cron.Start()
	// Heart Recursive execution
	//if ok {
	//	time.Sleep(10 * time.Second)
	//	go this.HeartChannel(ch, group, NodeChannel)
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
	//toDO: 删除15个区块前的列表
	if err != nil {
		defer utils.DeferRecoverLog("controller => time_task => createBlock:", err.Error(), 20001, nil)
		panic(err.Error())
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
		currentNode: node,
	}
	group.Done()
}

func (this *Controller) BlocksListen(from int64, to int64, api *blockchain.Api, node dao.InstanceNodes) error {
	var err error
	// 重新同步最近的12个区块
	if from > to {
		for i := to; i <= from; i++ {
			//fmt.Printf("Resync Block Create: %+v\n", i)
			this.SyncBlock(api, i, node)
		}
	} else {
		this.SyncBlock(api, from, node)
	}

	return err
}
