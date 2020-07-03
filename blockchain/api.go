package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

type Node struct {
	Address   string //地址
	Port      int    //端口
	ChainId   uint
	IsHttps   bool   //是否使用https
	NetworkId uint64 //链网络id
}

type Api struct {
	client       *rpc.Client
	address      string
	port         int
	chainId      uint
	simpleClient *ethclient.Client
	networkId    uint64
}

func NewApi(node *Node) (*Api, error) {
	var urlStr string
	if node.IsHttps {
		urlStr = fmt.Sprintf("https://%s:%d", node.Address, node.Port)
	} else {
		urlStr = fmt.Sprintf("http://%s:%d", node.Address, node.Port)
	}
	fmt.Println("now url=", urlStr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	client, err := rpc.DialContext(ctx, urlStr)
	if err != nil {
		return nil, err
	}
	return &Api{
		client:       client,
		address:      node.Address,
		port:         node.Port,
		chainId:      node.ChainId,
		networkId:    node.NetworkId,
		simpleClient: ethclient.NewClient(client),
	}, nil
}

func (this *Api) BlockByNumber(number *big.Int) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	block, err := this.simpleClient.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return block, nil
}
func (this *Api) TransactionByHash(hash common.Hash) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	transaction, _, err := this.simpleClient.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
func (this *Api) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	receipt, err := this.simpleClient.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
func (this *Api) ChangeNode(node *Node) error {
	if this.address == node.Address && this.port == node.Port {
		return nil
	}
	var urlStr string
	if node.IsHttps {
		urlStr = fmt.Sprintf("https://%s:%d", node.Address, node.Port)
	} else {
		urlStr = fmt.Sprintf("http://%s:%d", node.Address, node.Port)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	client, err := rpc.DialContext(ctx, urlStr)
	if err != nil {
		return err
	}
	this.client = client
	this.address = node.Address
	this.port = node.Port
	this.chainId = node.ChainId
	this.networkId = node.NetworkId
	this.simpleClient = ethclient.NewClient(client)
	return nil
}
func (this *Api) Close() {
	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
}
func (this *Api) GetNetworkId() uint64 {
	return this.networkId
}
func (this *Api) GetChainId() uint {
	return this.chainId
}
func (this *Api) GetPastEvents(query ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	log, err := this.simpleClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func (this *Api) GetHeaderByNumber() (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	header, err := this.simpleClient.HeaderByNumber(ctx, nil)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return header, nil
}
