package controllers

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"sipemanager/blockchain"
	"sipemanager/dao"
	"time"
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

func (this *Controller) createCrossEvent(nodes []dao.InstanceNodes) {
	a := time.Now()
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("current nodes %+v ", nodes[i])
		addresses := []common.Address{
			common.HexToAddress(nodes[i].CrossAddress),
		}
		node := ethereum.FilterQuery{
			FromBlock: big.NewInt(1),
			Addresses: addresses,
		}
		api, err := GetRpcApi(nodes[i])
		if err != nil {
			errors.New("cant not found nodes")
		}
		log, err := api.GetPastEvents(node)
		fmtLogs(log)
		fmt.Println("api ", api.GetChainId())

	}
	fmt.Println(time.Since(a))
}

func fmtLogs(logs []types.Log) {
	for _, log := range logs {
		fmt.Printf("%+v\n", log.BlockNumber)
	}
}
