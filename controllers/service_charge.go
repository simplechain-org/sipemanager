package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
	"github.com/gin-gonic/gin"
)

type ServiceChargeFee struct {
	AccumulatedFee *big.Int `json:"accumulated_fee"` //累计消耗手续费
	ReimbursedFee  *big.Int `json:"reimbursed_fee"`  //已报销
	CurrentFee     *big.Int `json:"current_fee"`     //本期应报销手续费
}

type ServiceChargeView struct {
	AnchorNodeId    uint   `json:"anchor_node_id"`   //锚定节点编号
	AnchorNodeName  string `json:"anchor_node_name"` //锚定节点名称，冗余方便查询
	TransactionHash string `json:"transaction_hash"` //交易哈希
	Fee             string `json:"fee"`              //报销手续费
	Coin            string `json:"coin"`             //报销的币种
	Sender          string `json:"sender"`           //出账账户地址
	Status          uint   `json:"status"`
	CreatedAt       string `json:"created_at"`
}

type ServiceChargeResult struct {
	TotalCount  int                 `json:"total_count"`  //总记录数
	CurrentPage int                 `json:"current_page"` //当前页数
	PageSize    int                 `json:"page_size"`    //页的大小
	PageData    []ServiceChargeView `json:"page_data"`    //页的数据
}


// @Summary 手续费报销记录
// @Tags ListServiceCharge
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query string true "锚定节点id"
// @Param current_page query string true "当前页"
// @Success 200 {object} JSONResult{data=ServiceChargeResult}
// @Router /service/charge/list [get]
func (this *Controller) ListServiceCharge(c *gin.Context) {
	//已报销手续费列表
	//报销时间
	//锚定节点名称
	//交易哈希
	//报销手续费
	//分页显示
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
	objects, err := this.dao.GetServiceChargeLogPage(start, pageSize, anchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	result := make([]ServiceChargeView, 0, len(objects))
	for _, obj := range objects {
		result = append(result, ServiceChargeView{
			AnchorNodeId:    obj.AnchorNodeId,
			AnchorNodeName:  obj.AnchorNodeName,
			TransactionHash: obj.TransactionHash,
			Fee:             obj.Fee,
			Coin:            obj.Coin,
			Sender:          obj.Sender,
			Status:          obj.Status,
			CreatedAt:       obj.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	count, err := this.dao.GetServiceChargeLogCount(anchorNodeId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	serviceChargeResult := &ServiceChargeResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    result,
	}
	this.echoResult(c, serviceChargeResult)
}

type AddServiceChargeParam struct {
	AnchorNodeId uint     `json:"anchor_node_id" form:"anchor_node_id"`
	NodeId       uint     `json:"node_id" form:"node_id"`
	WalletId     uint     `json:"wallet_id" form:"wallet_id"`
	Password     string   `json:"password" form:"password"`
	Fee          string `json:"fee" form:"fee"`   //报销手续费
	Coin         string   `json:"coin" form:"coin"` //报销的币种
}

//这里应注意到：一条链只能报销一种token(币)
// @Summary 新增手续费报销
// @Tags AddServiceCharge
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id formData uint true "锚定节点id"
// @Param node_id formData uint true "节点id"
// @Param wallet_id formData uint true "钱包id"
// @Param password formData string true "钱包密码"
// @Param fee formData string true "手续费"
// @Param coin formData string true "报销币种"
// @Success 200 {object} JSONResult{data=nil}
// @Router /service/charge/add [post]
func (this *Controller) AddServiceCharge(c *gin.Context) {
	var param AddServiceChargeParam
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
		fmt.Println("GetApiByNodeId:", err.Error())
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
		fmt.Println("GetWallet:", err.Error())
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		fmt.Println("GetPrivateKey:", err.Error())
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
	fee,success:=big.NewInt(0).SetString(param.Fee,10)
	if !success {
		this.echoError(c, errors.New("fee数据非法"))
		return
	}
	hash, err := source.AccumulateRewards(config, callerConfig, fee)
	if err != nil {
		this.echoError(c, err)
		return
	}
	serviceChargeLog := &dao.ServiceChargeLog{
		AnchorNodeId:    param.AnchorNodeId,
		AnchorNodeName:  anchorNode.Name,
		TransactionHash: hash,
		Fee:             param.Fee,
		Coin:            param.Coin,
		Sender:          address.String(),
	}
	id, err := this.dao.CreateServiceChargeLog(serviceChargeLog)
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
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := source.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateServiceChargeLogSourceStatus(id, uint(receipt.Status))
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
	}(source, id, hash)
	this.echoSuccess(c, "Success")
}

// @Summary 累计消耗手续费和累计已报销手续费及计算本期应报销手续费
// @Tags GetServiceChargeFee
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param anchor_node_id query uint true "锚定节点id"
// @Param node_id formData uint true "节点id"
// @Param coin formData string true "币种"
// @Success 200 {object} JSONResult{data=ServiceChargeFee}
// @Router /service/charge/fee [get]
func (this *Controller) GetServiceChargeFee(c *gin.Context) {

	coin := c.Query("coin")
	if coin == "" {
		this.echoError(c, errors.New("coin必须提供值"))
		return
	}
	anchorNodeIdStr := c.Query("anchor_node_id")

	if anchorNodeIdStr == "" {
		this.echoError(c, errors.New("anchor_node_id必须提供值"))
		return
	}
	anchorNodeId, err := strconv.ParseUint(anchorNodeIdStr, 10, 62)
	if err != nil {
		this.echoError(c, err)
		return
	}
	nodeIdStr := c.Query("node_id")

	if nodeIdStr == "" {
		this.echoError(c, errors.New("node_id必须提供值"))
		return
	}
	nodeId, err := strconv.ParseUint(nodeIdStr, 10, 62)
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetNodeById(uint(nodeId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	contract, err := this.dao.GetContractByChainId(node.ChainId)
	if err != nil {
		this.echoError(c, errors.New(fmt.Sprintf("获取链的合约失败 chain_id=%d", node.ChainId)))
		return
	}
	//获取锚定节点
	anchorNode, err := this.dao.GetAnchorNode(uint(anchorNodeId))
	if err != nil {
		this.echoError(c, errors.New(fmt.Sprintf("获取锚定节点信息失败 chain_id=%d", anchorNodeId)))
		return
	}
	//统计累计消耗手续费（从同步交易中获取）
	accumulatedFee, err := this.dao.GetTransactionSumFee(anchorNode.Address, contract.Address, "makerFinish", node.ChainId)
	if err != nil {
		this.echoError(c, errors.New("统计累计消耗手续费发生错误:"+err.Error()))
		return
	}
	//统计累计已报销手续费 (从报销表中获取，求和)
	reimbursedFee, err := this.dao.GetServiceChargeSumFee(uint(anchorNodeId), coin)
	if err != nil {
		this.echoError(c, err)
		return
	}
	result := &ServiceChargeFee{
		AccumulatedFee: accumulatedFee,                 //累计消耗手续费
		ReimbursedFee:  reimbursedFee,                  //累计已报销手续费
		CurrentFee:     big.NewInt(0).Sub(accumulatedFee,reimbursedFee), //计算本期应报销手续费
	}

	this.echoResult(c, result)
}
