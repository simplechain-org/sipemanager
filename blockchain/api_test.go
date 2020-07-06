package blockchain

import (
	"context"
	"fmt"

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
