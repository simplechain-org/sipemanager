package blockchain

import (
	"context"
	"fmt"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/simplechain-org/go-simplechain/ethclient"
	"github.com/simplechain-org/go-simplechain/rpc"
	"testing"
	"time"
)

func TestNewApi(t *testing.T) {
	urlStr := "ws://localhost:8546"
	//urlStr:="http://localhost:8545"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	client, err := rpc.DialContext(ctx, urlStr)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	ethClient := ethclient.NewClient(client)
	headerChan := make(chan *types.Header, 100)
	sub, err := ethClient.SubscribeNewHead(ctx, headerChan)
	if err != nil {
		t.Fatal(err)
	}
	for header := range headerChan {
		fmt.Println("header number", header.Number.String())
	}
	fmt.Println("sub=", sub)
}

func TestApi_GetMonitor(t *testing.T) {
	n := &Node{
		Address:   "192.168.4.107",
		Port:      8548,
		ChainId:   1,
		IsHttps:   false,
		NetworkId: 7,
	}
	api, err := NewApi(n)
	if err != nil {
		fmt.Println("err=", err)
	}
	result, err := api.GetMonitor()
	fmt.Println("result=", result)
}

func TestApi_LatestBalanceAt(t *testing.T) {
	n := &Node{
		Address:   "192.168.4.107",
		Port:      8548,
		ChainId:   1,
		IsHttps:   false,
		NetworkId: 7,
	}
	api, err := NewApi(n)
	if err != nil {
		fmt.Println("err=", err)
	}
	result, err := api.LatestBalanceAt(common.HexToAddress("0x17529b05513e5595ceff7f4fb1e06512c271a540"))
	fmt.Println("result=", result)
}
