package service

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"sipemanager/blockchain"
	"sipemanager/dao"
	"time"
)

type Service struct {
	Rpc *blockchain.Api
	dao *dao.DataBaseAccessObject
}

func ListenEvent(dao *dao.DataBaseAccessObject) {

	cron := cron.New()
	cron.AddFunc("@every 5s", func() {
		fmt.Println("current time is ", time.Now())
		nodes, err := dao.ListAllNode()
		if err != nil {
			errors.New("cant not found nodes")
		}
		go createCrossEvent(nodes)
	})
	cron.Start()
}

func GetRpcApi(node dao.Node) (*blockchain.Api, error) {
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

func createCrossEvent(nodes []dao.Node) {
	a := time.Now()
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("%+v\n", nodes[i].Address)
		api, err := GetRpcApi(nodes[i])
		if err != nil {
			errors.New("cant not found nodes")
		}
		fmt.Println("api ", api.GetChainId())
	}
	fmt.Println(time.Since(a))
}
