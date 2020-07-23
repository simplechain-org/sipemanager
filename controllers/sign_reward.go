package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
	"strings"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/gin-gonic/gin"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
)

type SignRewardView struct {
	AnchorNodeId    uint   `json:"anchor_node_id"`   //锚定节点编号
	AnchorNodeName  string `json:"anchor_node_name"` //锚定节点名称，冗余方便查询
	TransactionHash string `json:"transaction_hash"` //交易哈希
	TotalReward     string `gorm:"total_reward"`     //奖励池总额
	Rate            string `gorm:"rate"`             //签名量占比,存一个格式化后的结果
	Reward          string `gorm:"reward"`           //奖励值
	Coin            string `json:"coin"`             //奖励的币种
	Sender          string `json:"sender"`           //出账账户地址
	Status          uint   `json:"status"`
	CreatedAt       string `json:"created_at"`
}

type SignRewardResult struct {
	TotalCount  int                      `json:"total_count"`  //总记录数
	CurrentPage int                      `json:"current_page"` //当前页数
	PageSize    int                      `json:"page_size"`    //页的大小
	PageData    []*dao.SignRewardLogView `json:"page_data"`    //页的数据
}

// @Summary 签名奖励
// @Tags ListSignReward
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id formData string true "锚定节点id"
// @Param current_page formData string true "当前页"
// @Param page_size formData string true "页的记录数"
// @Success 200 {object} JsonResult{data=object}
// @Router /reward/list [get]
func (this *Controller) ListSignReward(c *gin.Context) {
	//分页显示
	//发放时间
	//锚定节点名称
	//奖励池总额度
	//签名量占比
	//奖励值
	//交易哈希
	var anchorNodeId uint
	//获取所有锚定节点的数据时，anchor_node_id设置为0
	anchorNodeIdStr := c.Query("anchor_node_id")
	if anchorNodeIdStr == "" {
		anchorNodeId = 0
	} else {
		id, err := strconv.ParseUint(anchorNodeIdStr, 10, 64)
		if err == nil {
			anchorNodeId = uint(id)
		}
	}
	var pageSize int = 10
	pageSizeStr := c.Query("page_size")
	if pageSizeStr != "" {
		size, err := strconv.ParseUint(pageSizeStr, 10, 64)
		if err == nil {
			pageSize = int(size)
			if pageSize > 100 {
				pageSize = 100
			}
		}
	}
	//当前页（默认为第一页）
	var currentPage int = 1
	currentPageStr := c.Query("current_page")
	if currentPageStr != "" {
		page, err := strconv.ParseUint(currentPageStr, 10, 64)
		if err == nil {
			currentPage = int(page)
		}
	}
	start := (currentPage - 1) * pageSize
	objects, err := this.dao.GetSignRewardLogPage(start, pageSize, anchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	count, err := this.dao.GetServiceChargeLogCount(anchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	serviceChargeResult := &SignRewardResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, serviceChargeResult)
}

// @Summary 剩余奖池总额
// @Tags GetTotalReward
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query string true "锚定节点id"
// @Param node_id query string true "节点id"
// @Success 200 {object} JsonResult{data=int}
// @Router /reward/total [get]
func (this *Controller) GetTotalReward(c *gin.Context) {
	anchorNodeIdStr := c.Query("anchor_node_id")
	nodeIdStr := c.Query("node_id")
	if anchorNodeIdStr == "" || nodeIdStr == "" {
		this.echoError(c, errors.New("缺少anchor_node_id或node_id"))
		return
	}
	anchorNodeId, err := strconv.ParseUint(anchorNodeIdStr, 10, 64)
	if err != nil {
		this.echoError(c, errors.New("非法的anchor_node_id"))
	}
	nodeId, err := strconv.ParseUint(nodeIdStr, 10, 64)
	if err != nil {
		this.echoError(c, errors.New("非法的node_id"))
	}
	anchorNode, err := this.dao.GetAnchorNode(uint(anchorNodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetNodeById(uint(nodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(uint(nodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	//链的合约
	contract, err := this.dao.GetContractByChainId(node.ChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var targetChainId uint
	if anchorNode.SourceChainId == node.ChainId {
		targetChainId = anchorNode.TargetChainId
	} else if anchorNode.TargetChainId == node.ChainId {
		targetChainId = anchorNode.SourceChainId
	}
	//目标链
	chain, err := this.dao.GetChain(targetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	config := &blockchain.AnchorNodeRewardConfig{
		AbiData:         []byte(contract.Abi),
		ContractAddress: common.HexToAddress(contract.Address),
		TargetNetworkId: chain.NetworkId,
		AnchorAddress:   common.HexToAddress(anchorNode.Address),
	}
	callerConfig := &blockchain.CallerConfig{
		NetworkId: source.GetNetworkId(),
	}
	totalReward, err := source.GetTotalReward(config, callerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, totalReward)
}

// @Summary 单笔签名奖励
// @Tags GetChainReward
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query string true "锚定节点id"
// @Param node_id query string true "节点id"
// @Success 200 {object} JsonResult{data=int}
// @Router /reward/chain [get]
func (this *Controller) GetChainReward(c *gin.Context) {
	anchorNodeIdStr := c.Query("anchor_node_id")
	nodeIdStr := c.Query("node_id")
	if nodeIdStr == "" {
		this.echoError(c, errors.New("缺少node_id"))
		return
	}
	anchorNodeId, err := strconv.ParseUint(anchorNodeIdStr, 10, 64)
	if err != nil {
		this.echoError(c, errors.New("非法的anchor_node_id"))
	}
	nodeId, err := strconv.ParseUint(nodeIdStr, 10, 64)
	if err != nil {
		this.echoError(c, errors.New("非法的node_id"))
	}
	anchorNode, err := this.dao.GetAnchorNode(uint(anchorNodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetNodeById(uint(nodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(uint(nodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	//链的合约
	contract, err := this.dao.GetContractByChainId(node.ChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var targetChainId uint
	if anchorNode.SourceChainId == node.ChainId {
		targetChainId = anchorNode.TargetChainId
	} else if anchorNode.TargetChainId == node.ChainId {
		targetChainId = anchorNode.SourceChainId
	}
	//目标链
	chain, err := this.dao.GetChain(targetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	config := &blockchain.AnchorNodeRewardConfig{
		AbiData:         []byte(contract.Abi),
		ContractAddress: common.HexToAddress(contract.Address),
		TargetNetworkId: chain.NetworkId,
		AnchorAddress:   common.HexToAddress(anchorNode.Address),
	}
	callerConfig := &blockchain.CallerConfig{
		NetworkId: source.GetNetworkId(),
	}
	reward, err := source.GetChainReward(config, callerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, reward)
}

type AnchorWorkCount struct {
	Rate      string   `json:"rate"`
	SignCount *big.Int `json:"sign_count"`
}

// @Summary 本期总签名数及签名工作量占比
// @Tags GetAnchorWorkCount
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query string true "锚定节点id"
// @Param node_id query string true "节点id"
// @Success 200 {object} JsonResult{data=AnchorWorkCount}
// @Router /anchor/work/count [get]
func (this *Controller) GetAnchorWorkCount(c *gin.Context) {
	anchorNodeIdStr := c.Query("anchor_node_id")
	nodeIdStr := c.Query("node_id")
	if nodeIdStr == "" {
		this.echoError(c, errors.New("缺少node_id"))
		return
	}
	anchorNodeId, err := strconv.ParseUint(anchorNodeIdStr, 10, 64)
	if err != nil {
		this.echoError(c, errors.New("非法的anchor_node_id"))
	}
	nodeId, err := strconv.ParseUint(nodeIdStr, 10, 64)
	if err != nil {
		this.echoError(c, errors.New("非法的node_id"))
	}
	anchorNode, err := this.dao.GetAnchorNode(uint(anchorNodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetNodeById(uint(nodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(uint(nodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	//链的合约
	contract, err := this.dao.GetContractByChainId(node.ChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var targetChainId uint
	if anchorNode.SourceChainId == node.ChainId {
		targetChainId = anchorNode.TargetChainId
	} else if anchorNode.TargetChainId == node.ChainId {
		targetChainId = anchorNode.SourceChainId
	}
	//目标链
	chain, err := this.dao.GetChain(targetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	config := &blockchain.AnchorNodeRewardConfig{
		AbiData:         []byte(contract.Abi),
		ContractAddress: common.HexToAddress(contract.Address),
		TargetNetworkId: chain.NetworkId,
		AnchorAddress:   common.HexToAddress(anchorNode.Address),
	}
	callerConfig := &blockchain.CallerConfig{
		NetworkId: source.GetNetworkId(),
	}
	signCount, finishCount, err := source.GetAnchorWorkCount(config, callerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//本期总签名数
	count := big.NewInt(0)

	count = count.Add(count, signCount)

	count = count.Add(count, finishCount)

	numerator := big.NewInt(0)

	numerator = numerator.Add(numerator, count)

	chainRegister, err := this.dao.GetChainRegisterByChaiId(anchorNode.SourceChainId, anchorNode.TargetChainId)

	anchorAddresses := strings.Split(chainRegister.AnchorAddresses, ",")

	for _, addr := range anchorAddresses {

		if addr == anchorNode.Address {
			continue
		}
		config := &blockchain.AnchorNodeRewardConfig{
			AbiData:         []byte(contract.Abi),
			ContractAddress: common.HexToAddress(contract.Address),
			TargetNetworkId: chain.NetworkId,
			AnchorAddress:   common.HexToAddress(addr),
		}
		callerConfig := &blockchain.CallerConfig{
			NetworkId: source.GetNetworkId(),
		}
		signCount, finishCount, err := source.GetAnchorWorkCount(config, callerConfig)
		if err != nil {
			this.echoError(c, err)
			return
		}
		count = count.Add(count, signCount)

		count = count.Add(count, finishCount)
	}
	//签名工作量占比
	rate := float64(numerator.Uint64()) / float64(count.Uint64())
	rateStr := fmt.Sprintf("%0.2f%%", rate)
	anchorWorkCount := &AnchorWorkCount{
		SignCount: count,
		Rate:      rateStr,
	}
	this.echoResult(c, anchorWorkCount)
}

type AddSignRewardParam struct {
	AnchorNodeId uint   `json:"anchor_node_id" form:"anchor_node_id"`
	NodeId       uint   `json:"node_id" form:"node_id"`
	WalletId     uint   `json:"wallet_id" form:"wallet_id"`
	Password     string `json:"password" form:"password"`
	Reward       string `json:"reward" form:"reward"` //奖励值
	Coin         string `json:"coin" form:"coin"`     //奖励币种

}

// @Summary 新增奖励发放
// @Tags AddSignReward
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id formData uint true "锚定节点id"
// @Param node_id formData uint true "节点id"
// @Param wallet_id formData uint true "钱包id"
// @Param password formData string true "钱包密码"
// @Param reward formData string true "奖励金额"
// @Param coin formData string true "奖励币种"
// @Success 200 {object} JsonResult{data=object}
// @Router /reward/add [post]
func (this *Controller) AddSignReward(c *gin.Context) {
	var param AddSignRewardParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(param.AnchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetNodeById(param.NodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(param.NodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//链的合约
	contract, err := this.dao.GetContractByChainId(node.ChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var targetChainId uint
	if anchorNode.SourceChainId == node.ChainId {
		targetChainId = anchorNode.TargetChainId
	} else if anchorNode.TargetChainId == node.ChainId {
		targetChainId = anchorNode.SourceChainId
	}
	chain, err := this.dao.GetChain(targetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		this.echoError(c, err)
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	config := &blockchain.AnchorNodeRewardConfig{
		AbiData:         []byte(contract.Abi),
		ContractAddress: common.HexToAddress(contract.Address),
		TargetNetworkId: source.GetNetworkId(),
		AnchorAddress:   common.HexToAddress(anchorNode.Address),
	}
	callerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  chain.NetworkId,
	}
	bigReward, success := big.NewInt(0).SetString(param.Reward, 10)
	if !success {
		this.echoError(c, errors.New("reward数据非法"))
		return
	}
	hash, err := source.AccumulateRewards(config, callerConfig, bigReward)
	if err != nil {
		this.echoError(c, err)
		return
	}
	signRewardLog := &dao.SignRewardLog{
		AnchorNodeId:    param.AnchorNodeId,
		AnchorNodeName:  anchorNode.Name,
		TransactionHash: hash,
		Reward:          param.Reward,
		Coin:            param.Coin,
		Sender:          address.String(),
	}
	id, err := this.dao.CreateSignRewardLog(signRewardLog)
	if err != nil {
		this.echoError(c, err)
		return
	}
	go func(source *blockchain.Api, id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			//时间到，做一下检测
			receipt, err := source.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateSignRewardLogStatus(id, uint(receipt.Status))
				if err != nil {
					continue
				}
				break
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(source, id, hash)
	this.echoSuccess(c, "Success")
}

type ConfigureSignRewardParam struct {
	AnchorNodeId uint   `json:"anchor_node_id" form:"anchor_node_id"`
	NodeId       uint   `json:"node_id" form:"node_id"`
	WalletId     uint   `json:"wallet_id" form:"wallet_id"`
	Password     string `json:"password" form:"password"`
	Reward       string `json:"reward" form:"reward"` //奖励值
	Coin         string `json:"coin" form:"coin"`     //奖励币种
}

// @Summary 配置签名奖励
// @Tags ConfigureSignReward
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id formData uint true "锚定节点id"
// @Param node_id formData uint true "节点id"
// @Param wallet_id formData uint true "钱包id"
// @Param password formData string true "钱包密码"
// @Param reward formData string true "奖励金额"
// @Param coin formData string true "奖励币种"
// @Success 200 {object} JsonResult{data=object}
// @Router /reward/configure [post]
func (this *Controller) ConfigureSignReward(c *gin.Context) {
	var param ConfigureSignRewardParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(param.AnchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetNodeById(param.NodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(param.NodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//链的合约
	contract, err := this.dao.GetContractByChainId(node.ChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var targetChainId uint
	if anchorNode.SourceChainId == node.ChainId {
		targetChainId = anchorNode.TargetChainId
	} else if anchorNode.TargetChainId == node.ChainId {
		targetChainId = anchorNode.SourceChainId
	}
	chain, err := this.dao.GetChain(targetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		this.echoError(c, err)
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	config := &blockchain.AnchorNodeRewardConfig{
		AbiData:         []byte(contract.Abi),
		ContractAddress: common.HexToAddress(contract.Address),
		TargetNetworkId: chain.NetworkId,
		AnchorAddress:   common.HexToAddress(anchorNode.Address),
	}
	callerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  source.GetNetworkId(),
	}
	reward, success := big.NewInt(0).SetString(param.Reward, 10)
	if !success {
		this.echoError(c, errors.New("reward数据非法"))
		return
	}
	hash, err := source.SetReward(config, callerConfig, reward)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoSuccess(c, hash)
}


