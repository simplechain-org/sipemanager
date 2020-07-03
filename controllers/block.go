package controllers

import (
	"math/big"
	"net/http"
	"strconv"
	"sync"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
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

func (this *Controller) BlocksListen(from, to int64, group *sync.WaitGroup) error {
	var err error
	//for i := int64(from); i <= to; i++ {
	//	fmt.Printf("Block Create: %+v\n", i)
	//}
	group.Done()
	return err
}
