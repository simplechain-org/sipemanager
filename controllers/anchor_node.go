package controllers

import (
	"errors"
	"fmt"
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

const (
	ANCHOR_NODE_ID_NOT_EXISTS_ERROR          int = 17001 //锚定节点不存在
	ANCHOR_NODE_ID_NOT_IN_CHAIN_EXISTS_ERROR int = 17002 //指定的锚定节点不在所选的节点所在的链上
)

type AddAnchorNodeParam struct {
	//发起链
	SourceChainId uint `json:"source_chain_id" form:"source_chain_id"`
	//目标链
	TargetChainId uint `json:"target_chain_id" form:"target_chain_id"`
	//发起链的节点id
	SourceNodeId uint `json:"source_node_id" form:"source_node_id"`
	//目标链的节点id
	TargetNodeId uint `json:"target_node_id" form:"target_node_id"`
	//签名账户地址
	AnchorAddress string `json:"anchor_address" form:"anchor_address"`
	//签名账户名称
	AnchorName string `json:"anchor_name" form:"anchor_name"`
	//钱包id
	WalletId uint `json:"wallet_id" form:"wallet_id"`
	//钱包密码
	Password     string `json:"password" form:"password"`
	SourceRpcUrl string `json:"source_rpc_url" form:"source_rpc_url"`
	TargetRpcUrl string `json:"target_rpc_url" form:"target_rpc_url"`
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
// @Param source_rpc_url formData string true "源链RpcUrl"
// @Param target_rpc_url formData string true "目标链RpcUrl"
// @Success 200 {object} JsonResult{data=object}
// @Router /anchor/node/add [post]
func (this *Controller) AddAnchorNode(c *gin.Context) {
	var param AddAnchorNodeParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	param.AnchorAddress = strings.TrimSpace(param.AnchorAddress)
	if !common.IsHexAddress(param.AnchorAddress) {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("锚定节点地址不合法"))
		return
	}
	// 去除空格
	param.AnchorAddress = strings.Replace(param.AnchorAddress, " ", "", -1)
	// 去除换行符
	param.AnchorAddress = strings.Replace(param.AnchorAddress, "\n", "", -1)
	//调用合约增加锚定节点，要注意是双链
	//添加到数据库
	sourceContract, err := this.dao.GetContractByChainId(param.SourceChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_ID_NOT_EXISTS_ERROR, fmt.Errorf("获取发起链失败，chain_id=%d", param.SourceChainId))
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		this.ResponseError(c, WALLET_ID_NOT_EXISTS_ERROR, fmt.Errorf("钱包id=%d的记录不存在", param.WalletId))
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		this.ResponseError(c, WALLET_PASSWORD_ERROR, fmt.Errorf("钱包解锁失败:%s", err.Error()))
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	source, err := this.getApiByNodeId(param.SourceNodeId)
	if err != nil {
		this.ResponseError(c, CHAIN_ID_NOT_EXISTS_ERROR, fmt.Errorf("发起链的节点id=%d创建api失败:%s", param.SourceNodeId, err.Error()))
		return
	}
	target, err := this.getApiByNodeId(param.TargetNodeId)
	if err != nil {
		this.ResponseError(c, CHAIN_ID_NOT_EXISTS_ERROR, fmt.Errorf("目标链的节点id=%d创建api失败:%s", param.TargetNodeId, err.Error()))
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
		this.ResponseError(c, CONTRACT_INVOKE_ERROR, fmt.Errorf("调用发起链的合约添加锚定节点失败:%s", err.Error()))
		return
	}
	//目标链合约
	targetContract, err := this.dao.GetContractByChainId(param.TargetChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, fmt.Errorf("获取目标链的合约失败，请确认合约已经部署，并配置到链 chain_id=%d", param.TargetChainId))
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
		this.ResponseError(c, CONTRACT_INVOKE_ERROR, fmt.Errorf("调用目标链的合约添加锚定节点失败:%s", err.Error()))
		return
	}
	pledge, err := this.dao.GetAnchorNodePledge(param.SourceChainId, param.TargetChainId, sourceContract.Address)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	anchorNode := &dao.AnchorNode{
		Name:          param.AnchorName,
		Address:       param.AnchorAddress,
		SourceChainId: param.SourceChainId,
		TargetChainId: param.TargetChainId,
		SourceHash:    sourceHash,
		TargetHash:    targetHash,
		SourceRpcUrl:  param.SourceRpcUrl,
		TargetRpcUrl:  param.TargetRpcUrl,
		Status:        true,
		Pledge:        pledge,
	}
	id, err := this.dao.CreateAnchorNode(anchorNode)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, fmt.Errorf("保存锚定节点数据失败：%s", err.Error()))
		return
	}
	go func(api *blockchain.Api, id uint, hash string, address string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 30
		i := 0
		for {
			<-ticker.C
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateSourceStatus(id, uint(receipt.Status))
				if err != nil {
					fmt.Println(err)
					continue
				}
				anchorNode, err := this.dao.GetAnchorNode(id)
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = this.dao.UpdateChainRegisterAnchorAddresses(anchorNode.SourceChainId, anchorNode.TargetChainId, address, "+", id)
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
	}(source, id, sourceHash, sourceContract.Address)

	go func(api *blockchain.Api, id uint, hash string, address string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 30
		i := 0
		for {
			<-ticker.C
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateTargetStatus(id, uint(receipt.Status))
				if err != nil {
					fmt.Println(err)
					continue
				}
				anchorNode, err := this.dao.GetAnchorNode(id)
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = this.dao.UpdateChainRegisterAnchorAddresses(anchorNode.TargetChainId, anchorNode.SourceChainId, address, "+", id)
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
	}(target, id, targetHash, targetContract.Address)
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
		this.ResponseError(c, REQUEST_PARAM_ERROR, fmt.Errorf("参数类型和值不匹配:%s", err.Error()))
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(param.AnchorNodeId)
	if err != nil {
		this.ResponseError(c, ANCHOR_NODE_ID_NOT_EXISTS_ERROR, fmt.Errorf("anchor_node_id=%d的记录没有找到", param.AnchorNodeId))
		return
	}
	sourceContract, err := this.dao.GetContractByChainId(anchorNode.SourceChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, fmt.Errorf("id=%d的链的跨链合约没有找到，请先配置跨链合约", anchorNode.SourceChainId))
		return
	}
	targetContract, err := this.dao.GetContractByChainId(anchorNode.TargetChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, fmt.Errorf("id=%d的链的跨链合约没有找到，请先配置跨链合约", anchorNode.TargetChainId))
		return
	}
	source, err := this.getApiByNodeId(param.SourceNodeId)
	if err != nil {
		this.ResponseError(c, NODE_ID_EXISTS_ERROR, fmt.Errorf("node_id=%d创建api失败:%s", param.SourceNodeId, err.Error()))
		return
	}
	target, err := this.getApiByNodeId(param.TargetNodeId)
	if err != nil {
		this.ResponseError(c, NODE_ID_EXISTS_ERROR, fmt.Errorf("node_id=%d创建api失败:%s", param.TargetNodeId, err.Error()))
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		this.ResponseError(c, WALLET_ID_NOT_EXISTS_ERROR, fmt.Errorf("wallet_id=%d的钱包不存在", param.WalletId))
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		this.ResponseError(c, WALLET_PASSWORD_ERROR, fmt.Errorf("钱包解锁失败:%s", err.Error()))
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainRegister, err := this.dao.GetChainRegisterWithAddress(anchorNode.SourceChainId, anchorNode.TargetChainId, sourceContract.Address, 1)
	if err != nil {
		this.ResponseError(c, ANCHOR_NODE_ID_NOT_IN_CHAIN_EXISTS_ERROR, err)
		return
	}
	ids := strings.Split(chainRegister.AnchorAddresses, ",")
	anchorNodeIds := fmt.Sprintf("%d", anchorNode.ID)
	var exists bool
	for _, id := range ids {
		if id == anchorNodeIds {
			exists = true
		}
	}
	if !exists {
		this.ResponseError(c, ANCHOR_NODE_ID_NOT_IN_CHAIN_EXISTS_ERROR, errors.New("当前选中的锚定节点不是当前链的锚定节点"))
		return
	}
	chainRegister, err = this.dao.GetChainRegisterWithAddress(anchorNode.SourceChainId, anchorNode.TargetChainId, targetContract.Address, 1)
	if err != nil {
		this.ResponseError(c, ANCHOR_NODE_ID_NOT_IN_CHAIN_EXISTS_ERROR, err)
		return
	}
	ids = strings.Split(chainRegister.AnchorAddresses, ",")
	anchorNodeIds = fmt.Sprintf("%d", anchorNode.ID)
	for _, id := range ids {
		if id == anchorNodeIds {
			exists = true
		}
	}
	if !exists {
		this.ResponseError(c, ANCHOR_NODE_ID_NOT_IN_CHAIN_EXISTS_ERROR, errors.New("当前选中的锚定节点不是当前链的锚定节点"))
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
	sourceHash, err := source.RemoveAnchors(sourceConfig, sourceCallerConfig)
	if err != nil {
		this.ResponseError(c, CONTRACT_INVOKE_ERROR, fmt.Errorf("调用发起链的合约移除锚定节点失败：%s", err.Error()))
		return
	}
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
		this.ResponseError(c, CONTRACT_INVOKE_ERROR, fmt.Errorf("调用目标链的合约移除锚定节点失败：%s", err.Error()))
		return
	}
	go func(source *blockchain.Api, target *blockchain.Api, id uint, sourceHash string, targetHash string, sourceAddress string, targetAddress string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 30
		i := 0
		for {
			<-ticker.C
			receipt, err := source.TransactionReceipt(common.HexToHash(sourceHash))
			if err == nil && receipt != nil {
				if receipt.Status == 1 {
					receipt, err := target.TransactionReceipt(common.HexToHash(targetHash))
					if err == nil && receipt != nil {
						if receipt.Status == 1 {
							anchorNode, err := this.dao.GetAnchorNode(id)
							if err != nil {
								fmt.Println(err)
								continue
							}
							err = this.dao.UpdateChainRegisterAnchorAddresses(anchorNode.SourceChainId, anchorNode.TargetChainId, sourceAddress, "-", id)
							if err != nil {
								fmt.Println(err)
								continue
							}
							err = this.dao.UpdateChainRegisterAnchorAddresses(anchorNode.TargetChainId, anchorNode.SourceChainId, targetAddress, "-", id)
							if err != nil {
								fmt.Println(err)
								continue
							}
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
	}(source, target, param.AnchorNodeId, sourceHash, targetHash, sourceContract.Address, targetContract.Address)
	this.echoSuccess(c, "Success")
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

	ID uint `json:"id"`

	//归属链A
	ChainAId uint `json:"chain_a_id"`
	//归属链B
	ChainBId uint `json:"chain_b_id"`

	SourceRpcUrl string `json:"source_rpc_url"`
	TargetRpcUrl string `json:"target_rpc_url"`
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
// @Router /anchor/node/list [get]
func (this *Controller) ListAnchorNode(c *gin.Context) {
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

	objects, err := this.dao.GetAnchorNodePage(start, pageSize, anchorNodeId)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, fmt.Errorf("分页获取锚定节点数据失败:%s", err.Error()))
		return
	}
	count, err := this.dao.GetAnchorNodeCount(anchorNodeId)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, fmt.Errorf("获取锚定节点总数失败:%s", err.Error()))
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
			SourceRpcUrl:   obj.SourceRpcUrl,
			TargetRpcUrl:   obj.TargetRpcUrl,
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
// @Tags GetAnchorNode
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query string true "锚定节点id"
// @Success 200 {object} JsonResult{data=AnchorNodeInfo}
// @Router /anchor/node/obtain [get]
func (this *Controller) GetAnchorNode(c *gin.Context) {
	anchorNodeIdStr := c.Query("anchor_node_id")
	if anchorNodeIdStr == "" {
		this.ResponseError(c, REQUEST_PARAM_ERROR, errors.New("anchor_node_id数据非法"))
		return
	}
	id, err := strconv.ParseUint(anchorNodeIdStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, fmt.Errorf("anchor_node_id进行数据转换时出错:%s", err.Error()))
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(uint(id))
	if err != nil {
		this.ResponseError(c, ANCHOR_NODE_ID_NOT_EXISTS_ERROR, fmt.Errorf("id为%d对应的锚定节点不存在", id))
		return
	}
	chainA, err := this.dao.GetChain(anchorNode.SourceChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_ID_NOT_EXISTS_ERROR, fmt.Errorf("找不到锚定节点的发起链chain_id=%d", anchorNode.SourceChainId))
		return
	}
	chainB, err := this.dao.GetChain(anchorNode.TargetChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_ID_NOT_EXISTS_ERROR, fmt.Errorf("找不到锚定节点的目标链chain_id=%d", anchorNode.TargetChainId))
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
		this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, fmt.Errorf("找不到发起链的合约链chain_id=%d", anchorNode.SourceChainId))
		return
	}
	makeFinishA, err := this.dao.GetTransactionSumFee(anchorNode.Address, contractA.Address, "makerFinish", anchorNode.SourceChainId)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, fmt.Errorf("获取makerFinish交易数失败：%s", err.Error()))
		return
	}
	reimbursedFeeA, err := this.dao.GetServiceChargeSumFee(anchorNode.ID, chainA.Symbol)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, fmt.Errorf("计算总的手续费失败:%s", err.Error()))
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
		this.ResponseError(c, DATABASE_ERROR, fmt.Errorf("根据链id=%d获取可用的节点失败", anchorNode.SourceChainId))
		return
	}
	source, err := this.getApiByNodeId(sourceNode.ID)
	if err != nil {
		this.ResponseError(c, NODE_ID_EXISTS_ERROR, fmt.Errorf("根据节点id创建api失败:%s", err.Error()))
		return
	}

	callerConfig := &blockchain.CallerConfig{
		NetworkId: chainA.NetworkId,
	}
	signCount, finishCount, err := source.GetAnchorWorkCount(config, callerConfig)
	if err != nil {
		this.ResponseError(c, CONTRACT_INVOKE_ERROR, fmt.Errorf("查询签名数失败:%s", err.Error()))
		return
	}
	//本期总签名数
	count := big.NewInt(0)

	count = count.Add(count, signCount)

	count = count.Add(count, finishCount)

	rewardA, err := this.dao.GetSignRewardLogSumFee(anchorNode.ID, chainA.Symbol)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
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
		this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, err)
		return
	}
	makeFinishB, err := this.dao.GetTransactionSumFee(anchorNode.Address, contractB.Address, "makerFinish", anchorNode.TargetChainId)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}

	reimbursedFeeB, err := this.dao.GetServiceChargeSumFee(anchorNode.ID, chainB.Symbol)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
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
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	target, err := this.getApiByNodeId(targetNode.ID)
	if err != nil {
		this.ResponseError(c, NODE_ID_EXISTS_ERROR, err)
		return
	}
	callerConfig = &blockchain.CallerConfig{
		NetworkId: chainB.NetworkId,
	}
	signCount, finishCount, err = target.GetAnchorWorkCount(config, callerConfig)
	if err != nil {
		this.ResponseError(c, CONTRACT_INVOKE_ERROR, err)
		return
	}
	//本期总签名数
	count = big.NewInt(0)

	count = count.Add(count, signCount)

	count = count.Add(count, finishCount)

	rewardB, err := this.dao.GetSignRewardLogSumFee(anchorNode.ID, chainB.Symbol)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
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

type UpdateAnchorNodeParam struct {
	ID           uint   `json:"id" form:"id"`
	SourceRpcUrl string `json:"source_rpc_url" form:"source_rpc_url"`
	TargetRpcUrl string `json:"target_rpc_url" form:"target_rpc_url"`
}

//锚定节点编辑
// @Summary 锚定节点编辑(RpcUrl)
// @Tags AddAnchorNode
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id formData uint true "anchorNodeId"
// @Param source_rpc_url formData string true "源链RpcUrl"
// @Param target_rpc_url formData string true "目标链RpcUrl"
// @Success 200 {object} JsonResult{data=object}
// @Router /anchor/node/update [post]
func (this *Controller) UpdateAnchorNode(c *gin.Context) {
	var param UpdateAnchorNodeParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, fmt.Errorf("数据类型不匹配:%s", err.Error()))
		return
	}
	err := this.dao.UpdateAnchorNode(param.ID, param.SourceRpcUrl, param.TargetRpcUrl)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, fmt.Errorf("修改锚定节点的rpc url失败:%s", err.Error()))
		return
	}
	this.echoSuccess(c, "Success")
}

// @Summary 获取所有锚定节点
// @Tags ListAnchorNode
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=dao.AnchorNode}
// @Router /anchor/node/list/all [get]
func (this *Controller) ListAllAnchorNode(c *gin.Context) {
	anchorNodes, err := this.dao.ListAnchorNode()
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, anchorNodes)
}
