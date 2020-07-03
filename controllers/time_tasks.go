package controllers

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"math/big"
	"sipemanager/blockchain"
	"sipemanager/dao"
	"strings"
	"sync"
	"time"
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
