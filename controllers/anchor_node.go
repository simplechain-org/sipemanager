package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/gin-gonic/gin"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
)

type AddAnchorNodeParam struct {
	SourceChainId uint   `json:"source_chain_id" form:"source_chain_id"`
	TargetChainId uint   `json:"target_chain_id" form:"target_chain_id"`
	SourceNodeId  uint   `json:"source_node_id" form:"source_node_id"`
	TargetNodeId  uint   `json:"target_node_id" form:"target_node_id"`
	AnchorAddress string `json:"anchor_address" form:"anchor_address"`
	AnchorName    string `json:"anchor_name" form:"anchor_name"`
	WalletId      uint   `json:"wallet_id" form:"wallet_id"`
	Password      string `json:"password" form:"password"`
}

//锚定节点管理
// @Summary 新增锚定节点
// @Tags AddAnchorNode
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param source_chain_id formData uint true "源链id"
// @Param target_chain_id formData uint true "目标链id"
// @Param source_node_id formData uint true "源节点id"
// @Param target_node_id formData uint true "目标链节点id"
// @Param anchor_address formData string true "锚定地址"
// @Param anchor_name formData string true "锚定节点名称"
// @Param wallet_id formData uint true "钱包id"
// @Param password formData string true "钱包密码"
// @Success 200 {object} JsonResult{data=object}
// @Router /anchor/node/add [post]
func (this *Controller) AddAnchorNode(c *gin.Context) {
	var param AddAnchorNodeParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	//调用合约增加锚定节点，要注意是双链
	//添加到数据库
	sourceContract, err := this.dao.GetContractByChainId(param.SourceChainId)
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

	source, err := this.getApiByNodeId(param.SourceNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	target, err := this.getApiByNodeId(param.TargetNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	configSource := &blockchain.AnchorNodeConfig{
		AbiData:         []byte(sourceContract.Abi),
		ContractAddress: common.HexToAddress(sourceContract.Address),
		TargetNetworkId: target.GetNetworkId(),
		AnchorAddresses: []common.Address{common.HexToAddress(param.AnchorAddress)},
	}
	callerConfigSource := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  source.GetNetworkId(),
	}
	sourceHash, err := source.AddAnchors(configSource, callerConfigSource)
	if err != nil {
		this.echoError(c, err)
		return
	}
	targetContract, err := this.dao.GetContractByChainId(param.TargetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	configTarget := &blockchain.AnchorNodeConfig{
		AbiData:         []byte(targetContract.Abi),
		ContractAddress: common.HexToAddress(targetContract.Address),
		TargetNetworkId: source.GetNetworkId(),
		AnchorAddresses: []common.Address{common.HexToAddress(param.AnchorAddress)},
	}
	callerConfigTarget := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  target.GetNetworkId(),
	}
	targetHash, err := target.AddAnchors(configTarget, callerConfigTarget)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//todo 万一其中一条链添加失败，这个时候该怎么处理

	anchorNode := &dao.AnchorNode{
		Name:          param.AnchorName,
		Address:       param.AnchorAddress,
		SourceChainId: param.SourceChainId,
		TargetChainId: param.TargetChainId,
		SourceHash:    sourceHash,
		TargetHash:    targetHash,
	}
	id, err := this.dao.CreateAnchorNode(anchorNode)
	if err != nil {
		this.echoError(c, err)
		return
	}

	go func(api *blockchain.Api, id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateSourceStatus(id, uint(receipt.Status))
				if err != nil {
					fmt.Println(err)
					continue
				}
				break
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(source, id, sourceHash)

	go func(api *blockchain.Api, id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateTargetStatus(id, uint(receipt.Status))
				if err != nil {
					fmt.Println(err)
					continue
				}
				break
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(target, id, targetHash)
	this.echoSuccess(c, "Success")
}

type RemoveAnchorNodeParam struct {
	AnchorNodeId uint   `json:"anchor_node_id" form:"anchor_node_id"`
	SourceNodeId uint   `json:"source_node_id" form:"source_node_id"`
	TargetNodeId uint   `json:"target_node_id" form:"target_node_id"`
	WalletId     uint   `json:"wallet_id" form:"wallet_id"`
	Password     string `json:"password" form:"password"`
}

// 锚定节点管理
// @Summary 删除锚定节点
// @Tags RemoveAnchorNode
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id formData uint true "锚定节点id"
// @Param source_node_id formData uint true "源节点id"
// @Param target_node_id formData uint true "目标链节点id"
// @Param wallet_id formData uint true "钱包id"
// @Param password formData string true "钱包密码"
// @Success 200 {object} JsonResult{data=object}
// @Router /anchor/node/remove [post]
func (this *Controller) RemoveAnchorNode(c *gin.Context) {
	//调用合约增加锚定节点，要注意是双链
	//从数据库中删除
	var param RemoveAnchorNodeParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(param.AnchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	sourceContract, err := this.dao.GetContractByChainId(anchorNode.SourceChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	targetContract, err := this.dao.GetContractByChainId(anchorNode.TargetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(param.SourceNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	target, err := this.getApiByNodeId(param.TargetNodeId)
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

	targetConfig := &blockchain.AnchorNodeConfig{
		AbiData:         []byte(targetContract.Abi),
		ContractAddress: common.HexToAddress(targetContract.Address),
		TargetNetworkId: source.GetNetworkId(),
		AnchorAddresses: []common.Address{common.HexToAddress(anchorNode.Address)},
	}
	targetCallerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  target.GetNetworkId(),
	}
	targetHash, err := target.RemoveAnchors(targetConfig, targetCallerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	sourceConfig := &blockchain.AnchorNodeConfig{
		AbiData:         []byte(sourceContract.Abi),
		ContractAddress: common.HexToAddress(sourceContract.Address),
		TargetNetworkId: target.GetNetworkId(),
		AnchorAddresses: []common.Address{common.HexToAddress(anchorNode.Address)},
	}
	sourceCallerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  source.GetNetworkId(),
	}
	sourceHash, err := target.RemoveAnchors(sourceConfig, sourceCallerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	go func(source *blockchain.Api, target *blockchain.Api, id uint, sourceHash string, targetHash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := source.TransactionReceipt(common.HexToHash(sourceHash))
			if err == nil && receipt != nil {
				if receipt.Status == 1 {
					receipt, err := target.TransactionReceipt(common.HexToHash(targetHash))
					if err == nil && receipt != nil {
						if receipt.Status == 1 {
							//两条链都已经完成合约调用成功，那么就移除删除数据库中锚定节点的数据
							err = this.dao.RemoveAnchorNode(id)
							if err != nil {
								fmt.Println("RemoveAnchorNode error:", err)
								continue
							}
							break
						}
					}
				}
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(source, target, param.AnchorNodeId, sourceHash, targetHash)
}

type AnchorNodeView struct {
	//创建时间
	CreatedAt string `json:"created_at"`
	//锚定节点名称
	AnchorNodeName string `json:"anchor_node_name"`
	//归属链A
	ChainA string `json:"chain_a"`
	//归属链B
	ChainB string `json:"chain_b"`

	//质押金额
	Pledge string `json:"pledge"`
	//身份状态
	Status string `json:"status"`

	ID uint `json:"ID"`

	//归属链A
	ChainAId uint `json:"chain_a_id"`
	//归属链B
	ChainBId uint `json:"chain_b_id"`
}

type AnchorNodeResult struct {
	TotalCount  int              `json:"total_count"`  //总记录数
	CurrentPage int              `json:"current_page"` //当前页数
	PageSize    int              `json:"page_size"`    //页的大小
	PageData    []AnchorNodeView `json:"page_data"`    //页的数据
}

// @Summary 锚定节点列表
// @Tags ListAnchorNode
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页"
// @Param page_size query string true "页的记录数"
// @Success 200 {object} JsonResult{data=AnchorNodeResult}
// @Router /service/charge/list [get]
func (this *Controller) ListAnchorNode(c *gin.Context) {
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

	objects, err := this.dao.GetAnchorNodePage(start, pageSize)
	if err != nil {
		this.echoError(c, err)
		return
	}
	count, err := this.dao.GetAnchorNodeCount()
	if err != nil {
		this.echoError(c, err)
		return
	}
	result := make([]AnchorNodeView, 0, len(objects))

	//分页显示，每页10条记录
	for _, obj := range objects {
		chainA, err := this.dao.GetChain(obj.SourceChainId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		chainB, err := this.dao.GetChain(obj.TargetChainId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var status string
		if obj.Status {
			status = "有效"
		} else {
			status = "无效"
		}
		result = append(result, AnchorNodeView{
			ID:             obj.ID,
			AnchorNodeName: obj.Name,
			ChainA:         chainA.Name,
			ChainB:         chainB.Name,
			CreatedAt:      obj.CreatedAt.Format("2006-01-02 15:04:05"),
			Pledge:         obj.Pledge,
			Status:         status,
			ChainAId:       obj.SourceChainId,
			ChainBId:       obj.TargetChainId,
		})
	}
	anchorNodeResult := &AnchorNodeResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    result,
	}
	this.echoResult(c, anchorNodeResult)
}

type AnchorNodeInfo struct {
	ID uint `json:"ID"`
	//创建时间
	CreatedAt string `json:"created_at"`
	//锚定节点名称
	AnchorNodeName string `json:"anchor_node_name"`
	//归属链A
	ChainA string `json:"chain_a"`
	//归属链B
	ChainB string `json:"chain_b"`

	//质押金额
	Pledge string `json:"pledge"`
	//身份状态
	Status string `json:"status"`

	ChainAInfo *ChainFeeInfo `json:"chain_a_info"`

	ChainBInfo *ChainFeeInfo `json:"chain_b_info"`
}
type ChainFeeInfo struct {
	//makefinish手续费
	MakeFinish string `json:"make_finish"`
	//已报销手续费
	ReimbursedFee string `json:"reimbursed_fee"`
	//累计有效签名
	ValidSignature string `json:"valid_signature"`
	//累计发放奖励
	Reward string `json:"reward"`

	ChainName string `json:"chain_name"`
}

// @Summary 锚定节点详情
// @Tags ListAnchorNode
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query string true "锚定节点id"
// @Success 200 {object} JsonResult{data=AnchorNodeInfo}
// @Router /anchor/node/obtain [get]
func (this *Controller) GetAnchorNode(c *gin.Context) {
	anchorNodeIdStr := c.Query("anchor_node_id")
	if anchorNodeIdStr == "" {
		this.echoError(c, errors.New("anchor_node_id数据非法"))
		return
	}
	id, err := strconv.ParseUint(anchorNodeIdStr, 10, 64)
	if err != nil {
		this.echoError(c, err)
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(uint(id))
	if err != nil {
		this.echoError(c, err)
		return
	}
	chainA, err := this.dao.GetChain(anchorNode.SourceChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	chainB, err := this.dao.GetChain(anchorNode.TargetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var status string
	if anchorNode.Status {
		status = "有效"
	} else {
		status = "无效"
	}

	contractA, err := this.dao.GetContractByChainId(anchorNode.SourceChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	makeFinishA, err := this.dao.GetTransactionSumFee(anchorNode.Address, contractA.Address, "makerFinish", anchorNode.SourceChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}

	reimbursedFeeA, err := this.dao.GetServiceChargeSumFee(anchorNode.ID, chainA.Symbol)
	if err != nil {
		this.echoError(c, err)
		return
	}

	config := &blockchain.AnchorNodeRewardConfig{
		AbiData:         []byte(contractA.Abi),
		ContractAddress: common.HexToAddress(contractA.Address),
		TargetNetworkId: chainB.NetworkId,
		AnchorAddress:   common.HexToAddress(anchorNode.Address),
	}

	//根据链id选择可以节点
	sourceNode, err := this.dao.GetNodeByChainId(anchorNode.SourceChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	source, err := this.getApiByNodeId(sourceNode.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}

	callerConfig := &blockchain.CallerConfig{
		NetworkId: chainA.NetworkId,
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

	rewardA, err := this.dao.GetSignRewardLogSumFee(anchorNode.ID, chainA.Symbol)
	if err != nil {
		this.echoError(c, err)
		return
	}
	chainFeeInfoA := &ChainFeeInfo{
		MakeFinish:     makeFinishA.String(),
		ReimbursedFee:  reimbursedFeeA.String(),
		ValidSignature: count.String(),
		Reward:         rewardA.String(),
	}

	contractB, err := this.dao.GetContractByChainId(anchorNode.TargetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	makeFinishB, err := this.dao.GetTransactionSumFee(anchorNode.Address, contractB.Address, "makerFinish", anchorNode.TargetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}

	reimbursedFeeB, err := this.dao.GetServiceChargeSumFee(anchorNode.ID, chainB.Symbol)
	if err != nil {
		this.echoError(c, err)
		return
	}

	config = &blockchain.AnchorNodeRewardConfig{
		AbiData:         []byte(contractB.Abi),
		ContractAddress: common.HexToAddress(contractB.Address),
		TargetNetworkId: chainA.NetworkId,
		AnchorAddress:   common.HexToAddress(anchorNode.Address),
	}
	//根据链id选择可以节点
	targetNode, err := this.dao.GetNodeByChainId(anchorNode.TargetChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	target, err := this.getApiByNodeId(targetNode.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	callerConfig = &blockchain.CallerConfig{
		NetworkId: chainB.NetworkId,
	}
	signCount, finishCount, err = target.GetAnchorWorkCount(config, callerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//本期总签名数
	count = big.NewInt(0)

	count = count.Add(count, signCount)

	count = count.Add(count, finishCount)

	rewardB, err := this.dao.GetSignRewardLogSumFee(anchorNode.ID, chainB.Symbol)
	if err != nil {
		this.echoError(c, err)
		return
	}
	chainFeeInfoB := &ChainFeeInfo{
		MakeFinish:     makeFinishB.String(),
		ReimbursedFee:  reimbursedFeeB.String(),
		ValidSignature: count.String(),
		Reward:         rewardB.String(),
	}
	anchorNodeInfo := &AnchorNodeInfo{
		ID:             anchorNode.ID,
		CreatedAt:      anchorNode.CreatedAt.Format(dateFormat),
		AnchorNodeName: anchorNode.Name,
		ChainA:         chainA.Name,
		ChainB:         chainB.Name,
		Status:         status,
		Pledge:         anchorNode.Pledge,
		ChainAInfo:     chainFeeInfoA,
		ChainBInfo:     chainFeeInfoB,
	}
	this.echoResult(c, anchorNodeInfo)
}
