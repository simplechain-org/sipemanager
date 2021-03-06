package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
	"math/big"
	"sipemanager/blockchain"
	"sipemanager/dao"
	"strconv"
)

const (
	PUNISHMENT_SUSPEND_ERROR           = 15001 //锚定节点签名功能已被禁用
	PUNISHMENT_RECOVERY_ERROR          = 15002 //当前已经暂停的锚定节点才能恢复
	PUNISHMENT_SUSPEND_DUPLICATE_ERROR = 15003 //重复暂停同一个锚定节点
	PUNISHMENT_PLEDGE_ERROR            = 15004 //token扣减数量非法
)

type AddPunishmentParam struct {
	AnchorNodeId uint   `json:"anchor_node_id" form:"anchor_node_id"`
	NodeId       uint   `json:"node_id" form:"node_id"`
	WalletId     uint   `json:"wallet_id" form:"wallet_id"`
	Password     string `json:"password" form:"password"`
	Value        string `json:"value" form:"value"` //扣减数量
	Coin         string `json:"coin" form:"coin"`   //扣减的币种
	ManageType   string `json:"manage_type" form:"manage_type"`
}

// @Summary 新增惩罚记录
// @Tags AddPunishment
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id formData uint true "锚定节点id"
// @Param node_id formData uint true "节点id"
// @Param wallet_id formData uint true "钱包id"
// @Param password formData string true "钱包密码"
// @Param value formData string true "扣减数量"
// @Param coin formData string true "扣减的币种"
// @Param manage_type formData string true "管理类型"
// @Success 200 {object} JsonResult{data=int}
// @Router /punishment/add [post]
func (this *Controller) AddPunishment(c *gin.Context) {
	var param AddPunishmentParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	anchorNode, err := this.dao.GetAnchorNode(param.AnchorNodeId)
	if err != nil {
		this.ResponseError(c, ANCHOR_NODE_ID_NOT_EXISTS_ERROR, err)
		return
	}
	node, err := this.dao.GetNodeById(param.NodeId)
	if err != nil {
		this.ResponseError(c, NODE_ID_EXISTS_ERROR, err)
		return
	}
	source, err := this.getApiByNodeId(param.NodeId)
	if err != nil {
		this.ResponseError(c, NODE_ID_EXISTS_ERROR, err)
		return
	}
	//suspend recovery
	if param.ManageType == "suspend" || param.ManageType == "recovery" {
		var status bool
		if param.ManageType == "recovery" {
			status = true
			if this.dao.PunishmentRecordNotFound(param.AnchorNodeId, "suspend") {
				this.ResponseError(c, PUNISHMENT_RECOVERY_ERROR, errors.New("当前已经暂停的锚定节点才能恢复"))
				return
			}
		}
		if param.ManageType == "suspend" {
			status = false
			if !this.dao.PunishmentRecordNotFound(param.AnchorNodeId, "suspend") {
				this.ResponseError(c, PUNISHMENT_SUSPEND_DUPLICATE_ERROR, errors.New("重复暂停同一个锚定节点"))
				return
			}
		}
		//链的合约
		contract, err := this.dao.GetContractByChainId(node.ChainId)
		if err != nil {
			this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, err)
			return
		}
		var targetChainId uint
		if anchorNode.SourceChainId == node.ChainId {
			targetChainId = anchorNode.TargetChainId
		} else if anchorNode.TargetChainId == node.ChainId {
			targetChainId = anchorNode.SourceChainId
		}
		sourceChain, err := this.dao.GetChain(node.ChainId)
		if err != nil {
			this.ResponseError(c, CHAIN_ID_NOT_EXISTS_ERROR, err)
			return
		}
		targetChain, err := this.dao.GetChain(targetChainId)
		if err != nil {
			this.ResponseError(c, CHAIN_ID_NOT_EXISTS_ERROR, err)
			return
		}
		wallet, err := this.dao.GetWallet(param.WalletId)
		if err != nil {
			this.ResponseError(c, WALLET_ID_NOT_EXISTS_ERROR, err)
			return
		}
		privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
		if err != nil {
			this.ResponseError(c, WALLET_PASSWORD_ERROR, err)
			return
		}
		address := crypto.PubkeyToAddress(privateKey.PublicKey)

		config := &blockchain.AnchorNodeRewardConfig{
			AbiData:         []byte(contract.Abi),
			ContractAddress: common.HexToAddress(contract.Address),
			TargetNetworkId: targetChain.NetworkId,
			AnchorAddress:   common.HexToAddress(anchorNode.Address),
		}
		callerConfig := &blockchain.CallerConfig{
			From:       address,
			PrivateKey: privateKey,
			NetworkId:  source.GetNetworkId(),
		}
		_, err = source.SetAnchorStatus(config, callerConfig, status)
		if err != nil {
			this.ResponseError(c, CONTRACT_INVOKE_ERROR, err)
			return
		}
		//对面链的合约
		targetContract, err := this.dao.GetContractByChainId(targetChainId)
		if err != nil {
			this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, err)
			return
		}

		node, err := this.dao.GetNodeByChainId(targetChainId)
		if err != nil {
			this.ResponseError(c, NODE_ID_EXISTS_ERROR, err)
			return
		}
		target, err := this.getApiByNodeId(node.ID)
		if err != nil {
			this.ResponseError(c, NODE_ID_EXISTS_ERROR, err)
			return
		}
		targetConfig := &blockchain.AnchorNodeRewardConfig{
			AbiData:         []byte(targetContract.Abi),
			ContractAddress: common.HexToAddress(targetContract.Address),
			TargetNetworkId: sourceChain.NetworkId,
			AnchorAddress:   common.HexToAddress(anchorNode.Address),
		}
		targetCallerConfig := &blockchain.CallerConfig{
			From:       address,
			PrivateKey: privateKey,
			NetworkId:  target.GetNetworkId(),
		}
		_, err = target.SetAnchorStatus(targetConfig, targetCallerConfig, status)
		if err != nil {
			this.ResponseError(c, CONTRACT_INVOKE_ERROR, err)
			return
		}
		if param.ManageType == "recovery" {
			err := this.dao.RemovePunishmentByManageType(param.AnchorNodeId, "suspend")
			if err != nil {
				this.ResponseError(c, DATABASE_ERROR, err)
				return
			}
		}
	}
	if param.ManageType == "token" {
		_, success := big.NewInt(0).SetString(param.Value, 10)
		if !success {
			this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("value数据非法"))
			return
		}
		//扣减
		err := this.dao.SubPledge(param.AnchorNodeId, param.Value)
		if err != nil {
			if err.Error() == "扣减数量非法" {
				this.ResponseError(c, PUNISHMENT_PLEDGE_ERROR, err)
				return
			} else {
				this.ResponseError(c, DATABASE_ERROR, err)
				return
			}
		}
	}
	punishment := &dao.Punishment{
		AnchorNodeId: param.AnchorNodeId,
		ManageType:   param.ManageType,
		Value:        param.Value,
		Coin:         param.Coin,
	}
	id, err := this.dao.CreatePunishment(punishment)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, id)
}

type PunishmentViewResult struct {
	TotalCount  int                   `json:"total_count"`  //总记录数
	CurrentPage int                   `json:"current_page"` //当前页数
	PageSize    int                   `json:"page_size"`    //页的大小
	PageData    []*dao.PunishmentView `json:"page_data"`    //页的数据
}

// @Summary 惩罚记录
// @Tags ListPunishment
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query string true "锚定节点id"
// @Param current_page query string true "当前页"
// @Success 200 {object} JsonResult{data=PunishmentViewResult}
// @Router /punishment/list [get]
func (this *Controller) ListPunishment(c *gin.Context) {
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

	objects, err := this.dao.GetPunishmentPage(start, pageSize, anchorNodeId)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	count, err := this.dao.GetPunishmentCount(anchorNodeId)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	punishmentViewResult := &PunishmentViewResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, punishmentViewResult)
}
