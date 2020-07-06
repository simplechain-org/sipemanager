package controllers

import (
	"fmt"

	"math/big"
	"net/http"
	"strconv"
	"strings"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	userClient map[uint]*blockchain.Api
	dao        *dao.DataBaseAccessObject
}

//获取一个连接到节点的连接
func (this *Controller) getBlockChainApi(userId uint) (*blockchain.Api, error) {
	api, ok := this.userClient[userId]
	if !ok {
		node, err := this.dao.GetUserCurrentNode(userId)
		if err != nil {
			return nil, err
		}
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
		this.userClient[userId] = api
		return api, nil
	}
	return api, nil
}

type Block struct {
	ParentHash   common.Hash    `json:"parentHash"`
	UncleHash    common.Hash    `json:"sha3Uncles"`
	CoinBase     common.Address `json:"miner"`
	Difficulty   *big.Int       `json:"difficulty"`
	Number       *big.Int       `json:"number"`
	GasLimit     uint64         `json:"gasLimit"`
	GasUsed      uint64         `json:"gasUsed"`
	Time         uint64         `json:"timestamp"`
	Nonce        uint64         `json:"nonce"`
	Transactions int            `json:"transactions"`
}

//分页获取区块列表
// @Summary 分页获取区块列表
// @Tags block
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} JSONResult{data=[]Block}
// @Router /block/list [get]
func (this *Controller) GetPageBlock(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	api, err := this.getBlockChainApi(user.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//获取最新区块
	block, err := api.BlockByNumber(nil)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var pageSize uint64 = 10
	//总记录数,区块号是从0开始计算的，所以这里要加1
	var total uint64 = block.Number().Uint64() + 1
	//总页数
	var totalPage = total / pageSize
	//如果不能除尽，那么就需要加1
	if total%pageSize != 0 {
		totalPage++
	}
	//当前页（默认为第一页）
	var currentPage uint64 = 1
	currentPageStr := c.Query("currentPage")
	if currentPageStr != "" {
		page, err := strconv.ParseUint(currentPageStr, 10, 64)
		if err == nil {
			currentPage = page
		}
	}
	var start uint64 = 0
	if total >= (currentPage-1)*pageSize+1 {
		start = total - (currentPage-1)*pageSize - 1
	}
	var end uint64 = 0
	if start >= pageSize {
		end = start - pageSize + 1
	} else {
		end = 0
	}
	blocks := make([]Block, 0)
	for i := start; i >= end; i-- {
		block, err := api.BlockByNumber(big.NewInt(0).SetUint64(i))
		if err != nil {
			continue
		}
		blocks = append(blocks, Block{
			ParentHash:   block.ParentHash(),
			UncleHash:    block.UncleHash(),
			CoinBase:     block.Coinbase(),
			Difficulty:   block.Difficulty(),
			Number:       block.Number(),
			GasLimit:     block.GasLimit(),
			GasUsed:      block.GasUsed(),
			Time:         block.Time(),
			Nonce:        block.Nonce(),
			Transactions: len(block.Transactions()),
		})
		if i == 0 {
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"result": map[string]interface{}{"data": blocks, "total": total},
		"code":   0,
	})
}

//目前采取的策略是:用户切换节点时，把原先已有的连接关闭，根据新的node信息重新建立一个连接
//如果没有切换节点就一直保持用已有的连接
func (this *Controller) onChangeNode(userId uint) (*blockchain.Api, error) {
	if _, ok := this.userClient[userId]; ok {
		this.userClient[userId].Close()
		delete(this.userClient, userId)
	}
	node, err := this.dao.GetUserCurrentNode(userId)
	if err != nil {
		return nil, err
	}
	n := &blockchain.Node{
		Address:   node.Address,
		Port:      node.Port,
		ChainId:   node.ChainId,
		IsHttps:   node.IsHttps,
		NetworkId: node.NetworkId,
	}
	api, err := blockchain.NewApi(n)
	this.userClient[userId] = api
	if err != nil {
		return nil, err
	}
	return api, nil
}

func (this *Controller) SyncBlock(api *blockchain.Api, number int64, node dao.InstanceNodes) {
	chainId := node.ChainId
	fmt.Printf("----当前写入区块号:%+v, ----当前ChainId： %+v ————\n", number, chainId)
	block, err := api.BlockByNumber(big.NewInt(0).SetInt64(number))
	if err != nil {
		logrus.Warn(&ErrLogCode{message: "time_task => SyncBlock:", code: 20004, err: err.Error()})
	}
	blockRecord := dao.Block{
		ParentHash:   block.ParentHash().Hex(),
		UncleHash:    block.UncleHash().Hex(),
		CoinBase:     block.Coinbase().Hex(),
		Difficulty:   block.Difficulty().Int64(),
		Number:       block.Number().Int64(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		Time:         block.Time(),
		Nonce:        block.Nonce(),
		Transactions: len(block.Transactions()),
		BlockHash:    block.Hash().Hex(),
		ChainId:      chainId,
	}
	replaceErr := this.dao.BlockReplace(blockRecord)
	if replaceErr != nil {
		fmt.Printf("%+v\n", blockRecord)
		logrus.Warn(&ErrLogCode{message: "time_task => SyncBlock:", code: 20005, err: err.Error()})
	}

	if len(block.Transactions()) > 0 {
		for index, transaction := range block.Transactions() {
			netId := api.GetNetworkId()
			// 获取from地址
			msg, err := transaction.AsMessage(types.NewEIP155Signer(big.NewInt(0).SetUint64(netId)))
			if err != nil {
				logrus.Error(err.Error())
			}
			from := msg.From().Hex()
			receipt, err := api.TransactionReceipt(transaction.Hash())
			if err != nil {
				logrus.Error("TransactionReceipt:", err)
			}
			// hexutil.Encode(transaction.Data()),
			txRecord := dao.Transaction{
				BlockHash:        block.Hash().Hex(),
				BlockNumber:      block.Number().Int64(),
				Hash:             transaction.Hash().Hex(),
				From:             strings.ToLower(from),
				GasUsed:          block.GasUsed(),
				GasPrice:         transaction.GasPrice().String(),
				Input:            transaction.Data(),
				Nonce:            transaction.Nonce(),
				TransactionIndex: int64(index),
				Value:            transaction.Value().String(),
				Timestamp:        block.Time(),
				Status:           receipt.Status,
				ChainId:          chainId,
			}
			if transaction.To() != nil {
				txRecord.To = strings.ToLower(transaction.To().Hex())
			}
			txReplaceErr := this.dao.TxReplace(txRecord)
			if txReplaceErr != nil {
				logrus.Error("Transactions Create:", err.Error())
			}

		}
	}

	if len(block.Uncles()) > 0 {
		for _, uncle := range block.Uncles() {
			uncleRecord := dao.Uncle{
				ParentHash:  uncle.ParentHash.Hex(),
				UncleHash:   uncle.UncleHash.Hex(),
				CoinBase:    uncle.Coinbase.Hex(),
				Difficulty:  uncle.Difficulty.Int64(),
				Number:      uncle.Number.Int64(),
				GasLimit:    uncle.GasLimit,
				GasUsed:     uncle.GasUsed,
				Time:        uncle.Time,
				Nonce:       uncle.Nonce.Uint64(),
				ChainId:     chainId,
				BlockNumber: block.Number().Int64(),
			}
			uncleReplaceErr := this.dao.UncleReplace(uncleRecord)
			if uncleReplaceErr != nil {
				logrus.Error("UncleRecord Create:", err.Error())
			}
		}

	}
}
