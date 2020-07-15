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
// @Success 200 {object} JSONResult{data=int}
// @Router /punishment/add [post]
func (this *Controller) AddPunishment(c *gin.Context) {
	var param AddPunishmentParam
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
	//suspend recovery
	if param.ManageType == "suspend" || param.ManageType == "recovery" {
		var status bool
		if param.ManageType == "recovery" {
			status = true
		}
		if param.ManageType == "suspend" {
			status = false
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
		_, err = source.SetAnchorStatus(config, callerConfig, status)
		if err != nil {
			this.echoError(c, err)
			return
		}
	}
	if param.ManageType == "token" {
		_, success := big.NewInt(0).SetString(param.Value, 10)
		if !success {
			this.echoError(c, errors.New("value数据非法"))
			return
		}
		//扣减
		err := this.dao.SubPledge(param.AnchorNodeId, param.Value)
		if err != nil {
			this.echoError(c, err)
			return
		}
	}
	punishment := &dao.Punishment{
		AnchorNodeId:   param.AnchorNodeId,
		AnchorNodeName: anchorNode.Name,
		ManageType:     param.ManageType,
		Value:          param.Value,
		Coin:           param.Coin,
	}
	id, err := this.dao.CreatePunishment(punishment)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, id)
}

type PunishmentView struct {
	Value string `json:"value"` //惩罚数量
	Coin  string `json:"coin"`  //惩罚币种
	//suspend recovery token
	ManageType     string `json:"manage_type"`      //管理类型
	AnchorNodeId   uint   `json:"anchor_node_id"`   //锚定节点编号
	AnchorNodeName string `json:"anchor_node_name"` //锚定节点名称，冗余方便查询
	CreatedAt      string `json:"created_at"`
}

type PunishmentViewResult struct {
	TotalCount  int                 `json:"total_count"`  //总记录数
	CurrentPage int                 `json:"current_page"` //当前页数
	PageSize    int                 `json:"page_size"`    //页的大小
	PageData    []PunishmentView `json:"page_data"`    //页的数据
}

// @Summary 惩罚记录
// @Tags ListPunishment
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query string true "锚定节点id"
// @Param current_page query string true "当前页"
// @Success 200 {object} JSONResult{data=PunishmentViewResult}
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
		this.echoError(c, err)
		return
	}
	result := make([]PunishmentView, 0, len(objects))
	for _, obj := range objects {
		result = append(result, PunishmentView{
			AnchorNodeId:   obj.AnchorNodeId,
			AnchorNodeName: obj.AnchorNodeName,
			Coin:           obj.Coin,
			ManageType:     obj.ManageType,
			CreatedAt:      obj.CreatedAt.Format(dateFormat),
		})
	}
	count, err := this.dao.GetPunishmentCount(anchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	punishmentViewResult := &PunishmentViewResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    result,
	}
	this.echoResult(c, punishmentViewResult)
}
