package controllers

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func (this *Controller) createCrossEvent(nodes []dao.InstanceNodes) {
	a := time.Now()
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("current nodes %+v ", nodes[i])
		contract, err := this.dao.GetContractById(nodes[i].ContractId)
		blockNumber := this.dao.GetMaxBlockNumber(nodes[i].ChainId)
		fmt.Printf("cxc%+v\n", blockNumber)
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
		fmt.Println("api ", api.GetChainId())

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

func (this *Controller) createBlock(nodes []dao.InstanceNodes, group *sync.WaitGroup) {
	a := time.Now()
	for i := 0; i < len(nodes); i++ {
		//sync all instance blocks
		api, err := GetRpcApi(nodes[i])
		header, err := api.GetHeaderByNumber()
		dbMaxNum, err := this.dao.GetNewBlockNumber(nodes[i].ChainId)
		if err != nil {
			logrus.Warn(&ErrLogCode{message: "time_task => createBlock:", code: 20001, err: err.Error()})
		}
		newBlockNumber := header.Number.Int64()

		fmt.Printf("$-----%+v\n", dbMaxNum)
		fmt.Printf("$-----%+v\n", newBlockNumber)
		var i int64
		for i = dbMaxNum; i <= newBlockNumber; i++ {
			var to int64
			if i == newBlockNumber {
				to = newBlockNumber - 12
			} else {
				to = newBlockNumber
			}
			logrus.Info(i, "from:", i)
			logrus.Infof("to:%v", to)
			group.Add(1)
			go this.BlocksListen(i, to, group)
		}
		group.Wait()
	}
	fmt.Println(time.Since(a))
}

func (this *Controller) BlocksListen(from, to int64, group *sync.WaitGroup) error {
	var err error
	//for i := int64(from); i <= to; i++ {
	//	fmt.Printf("Block Create: %+v\n", i)
	//}
	group.Done()
	return err
}
