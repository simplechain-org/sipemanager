package blockchain

import (
	"context"
	"fmt"
	"github.com/simplechain-org/go-simplechain/common/hexutil"
	"math/big"
	"time"

	"github.com/simplechain-org/go-simplechain"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/simplechain-org/go-simplechain/ethclient"
	"github.com/simplechain-org/go-simplechain/rpc"
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

func NewDirectApi(url string) (*Api, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	client, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Api{
		client:       client,
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
func (this *Api) GetPastEvents(query simplechain.FilterQuery) ([]types.Log, error) {
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

type Monitor struct {
	Tally    map[common.Address]uint64
	Recently map[common.Address]uint32
}

func (this *Api) GetMonitor() (Monitor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	var result Monitor
	err := this.client.CallContext(ctx, &result, "cross_monitor")
	if err != nil {
		logrus.Error(err.Error())
		return result, err
	}
	return result, err
}

func (this *Api) LatestBalanceAt(account common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	var result hexutil.Big
	err := this.client.CallContext(ctx, &result, "eth_getBalance", account, "latest")
	if err != nil {
		logrus.Error(err.Error())
	}
	return (*big.Int)(&result), err
}
