package controllers

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"
	"sipemanager/utils"

	"github.com/robfig/cron/v3"
	"github.com/simplechain-org/go-simplechain"
	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/sirupsen/logrus"
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
		dbMaxNums, err := this.dao.GetAllMaxBlockNumber()
		if err != nil {
			logrus.Warn(utils.ErrLogCode{LogType: "controller => time_task => ListenAnchors:", Code: 20011, Message: err.Error(), Err: nil})
		}
		for _, dbMax := range dbMaxNums {
			if dbMax.Number > 15 {
				this.dao.Delete(dbMax.Number-15, dbMax.ChainId)
			}
		}
	})
	cron.Start()
}

func (this *Controller) ListenHeartChannel() {
	for {
		ch, ok := <-this.NodeChannel
		fmt.Println("node channel is", ch.ChainId, ch.BlockNumber)
		logrus.Infof("node channel is %+v, ok = %+v", ch, ok)
		if ok {
			time.Sleep(time.Duration(5) * time.Second)
			chain, err := this.dao.GetChain(ch.ChainId)
			if err != nil {
				logrus.Warn(utils.ErrLogCode{LogType: "controller => time_task => ListenHeartChannel:", Code: 20009, Message: err.Error(), Err: nil})
			}
			if chain.ContractInstanceId == ch.ContractInstanceId {
				go this.HeartChannel(ch)
			}
		}
	}

}

//func (this *Controller) ListenStopChannel() {
//	for {
//		ch, _ := <-this.CloseChannel
//		NodeCh, _ := <-this.NodeChannel
//		logrus.Infof("stop ---- is %+v", ch)
//		if ch.Status && NodeCh.ContractInstanceId == ch.ContractInstanceId && NodeCh.ChainId == ch.ChainId {
//			defer func() {
//				this.onceClose.Do(func() {
//					close(this.NodeChannel)
//				})
//				this.onceClose.Do(func() {
//					close(this.NodeChannel)
//					fmt.Println("send goroutine closed !")
//					this.group.Done()
//				})
//			}()
//		}
//	}
//}

func (this *Controller) ListenDirectBlock() {
	nodes, err := this.dao.GetInstancesJoinNode()
	filterNodes := utils.RemoveRepByLoop(nodes)
	if err != nil {
		logrus.Warn(utils.ErrLogCode{LogType: "controller => time_task => ListenBlock:", Code: 20007, Message: err.Error(), Err: nil})
	}
	go this.createBlock(filterNodes)
}

func (this *Controller) UpdateDirectBlock(chainId uint) {
	nodes, err := this.dao.GetInstancesJoinNode()
	filterNodes := utils.RemoveRepByLoop(nodes)
	updateChains := make([]dao.InstanceNodes, 0)
	for _, item := range filterNodes {
		if item.ChainId == chainId {
			updateChains = append(updateChains, item)
		}
	}
	if err != nil {
		logrus.Warn(utils.ErrLogCode{LogType: "controller => time_task => ListenBlock:", Code: 20007, Message: err.Error(), Err: nil})
	}
	go this.createBlock(updateChains)
}

func (this *Controller) createCrossEvent(nodes []dao.InstanceNodes) {
	for _, node := range nodes {
		contract, err := this.dao.GetContractById(node.ContractId)
		blockNumber := this.dao.GetMaxCrossNumber(node.ChainId)
		addresses := []common.Address{
			common.HexToAddress(node.CrossAddress),
		}

		records := simplechain.FilterQuery{
			FromBlock: big.NewInt(blockNumber),
			Addresses: addresses,
		}
		api, err := GetRpcApi(node)
		logs, err := api.GetPastEvents(records)
		if err != nil {
			logrus.Warn(utils.ErrLogCode{LogType: "controller => time_task => FilterLogs:", Code: 20014, Message: err.Error(), Err: nil})
		}

		if len(logs) > 0 {
			abiParsed, err := abi.JSON(strings.NewReader(contract.Abi))
			if err != nil {
				logrus.Warn(utils.ErrLogCode{LogType: "controller => time_task => createCrossEvent:", Code: 20013, Message: "Unable to parse ABi normally", Err: err})
			}
			this.EventLog(logs, abiParsed, node)
		}
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

func (this *Controller) HeartChannel(ch BlockChannel) {

	//ch, ok := <-this.NodeChannel
	//logrus.Infof("node HeartChannel is %+v, ok = %+v", ch, ok)
	//select {
	//case <-NodeChannel:
	//	fmt.Println("消费完成……………………")
	//	return
	//case <-time.After(time.Second * 5):
	//	fmt.Println("超时………………………")
	//	return
	//}

	current, err := this.dao.GetNodeById(ch.NodeId)
	if err != nil {
		defer utils.DeferRecoverLog("controller => time_task => HeartChannel:", err.Error(), 20005, nil)
		panic(err.Error())
	}
	chain, err := this.dao.GetChain(current.ChainId)
	if err != nil {
		logrus.Error(err)
		return
	}

	currents := []dao.InstanceNodes{
		dao.InstanceNodes{
			Address:            current.Address,
			Port:               current.Port,
			IsHttps:            current.IsHttps,
			NetworkId:          chain.NetworkId,
			Name:               current.Name,
			ChainId:            current.ChainId,
			ContractId:         ch.CurrentNode.ContractId,
			CrossAddress:       ch.CurrentNode.CrossAddress,
			NodeId:             current.ID,
			ContractInstanceId: ch.ContractInstanceId,
		},
	}
	go this.createBlock(currents)
}

func (this *Controller) createBlock(nodes []dao.InstanceNodes) {
	for _, node := range nodes {
		this.group.Add(1)
		go this.syncAllNodes(node)
	}
}

func (this *Controller) syncAllNodes(node dao.InstanceNodes) {
	api, err := GetRpcApi(node)
	chainId := node.ChainId
	header, err := api.GetHeaderByNumber()
	dbMaxNum := this.dao.GetMaxBlockNumber(chainId)
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

	this.NodeChannel <- BlockChannel{
		ChainId:            chainId,
		NodeId:             node.NodeId,
		BlockNumber:        dbMaxNum,
		CurrentNode:        node,
		ContractInstanceId: node.ContractInstanceId,
	}

	this.group.Done()
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
