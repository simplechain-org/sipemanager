package controllers

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
)

type RegisterChainTwoWayParam struct {
	SourceChainId    uint `json:"source_chain_id" form:"source_chain_id"`
	TargetChainId    uint `json:"target_chain_id" form:"target_chain_id"`
	SourceNodeId     uint `json:"source_node_id" form:"source_node_id"`
	TargetNodeId     uint `json:"target_node_id" form:"target_node_id"`
	SignConfirmCount uint `json:"sign_confirm_count" form:"sign_confirm_count"`
	//锚定节点地址
	AnchorAddresses []string `json:"anchor_addresses" form:"anchor_addresses"`
	//锚定节点名称
	AnchorNames []string `json:"anchor_names" form:"anchor_names"`
	WalletId    uint     `json:"wallet_id" form:"wallet_id"`
	Password    string   `json:"password" form:"password"`
	Pledge      string   `json:"pledge"` //质押sipc的金额
}

//@Summary 注册新的跨链对
//@Tags RegisterChainTwoWay
//@Accept application/x-www-form-urlencoded
//@Accept application/json
//@Produce application/json
//@Param source_chain_id formData uint true "源链Id"
//@Param target_chain_id formData uint true "目标链Id"
//@Param source_node_id formData uint true "源链节点Id"
//@Param target_node_id formData uint true "目标链节点Id"
//@Param sign_confirm_count formData uint true "最小确认数"
//@Param anchor_addresses formData array true "锚定节点地址 字符串数组"
//@Param anchor_names formData array true "锚定节点名称 字符串数组"
//@Param pledge formData string true "锚定节点质押金额"
//@Param wallet_id formData uint true "钱包id"
//@Param password formData string true "钱包密码"
//@Success 200 {object} JsonResult
//@Router /contract/register/once [post]
func (this *Controller) RegisterChainTwoWay(c *gin.Context) {
	var param RegisterChainTwoWayParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
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

	errChan := make(chan error, 2)
	source, err := this.getApiByNodeId(param.SourceNodeId)
	if err != nil {
		this.ResponseError(c, NODE_ID_EXISTS_ERROR, err)
		return
	}
	target, err := this.getApiByNodeId(param.TargetNodeId)
	if err != nil {
		this.ResponseError(c, NODE_ID_EXISTS_ERROR, err)
		return
	}
	db := this.dao.BeginTransaction()
	ids := make([]string, 0)
	for index, address := range param.AnchorAddresses {
		anchorNode := &dao.AnchorNode{
			Name:          param.AnchorNames[index],
			Address:       address,
			SourceChainId: param.SourceChainId,
			TargetChainId: param.TargetChainId,
			Pledge:        param.Pledge,
			Status:        true, //默认为有效
		}
		id, err := this.dao.CreateAnchorNodeByTx(db, anchorNode)
		if err != nil {
			db.Rollback()
			this.ResponseError(c, DATABASE_ERROR, err)
			return
		}
		ids = append(ids, fmt.Sprintf("%d", id))
	}
	idString := strings.Join(ids, ",")
	pledge, ok := new(big.Int).SetString(param.Pledge, 10)
	if !ok {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("pledge 值非法"))
		return
	}
	sourceContract, err := this.dao.GetContractByChainId(param.SourceChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, err)
		return
	}
	targetContract, err := this.dao.GetContractByChainId(param.TargetChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_CONTRACT_NOT_EXISTS_ERROR, err)
		return
	}
	if !this.dao.ChainRegisterRecordNotFound(param.SourceChainId, param.TargetChainId, sourceContract.Address, 1) {
		this.ResponseError(c, CHAIN_REGISTER_EXISTS_ERROR, errors.New("链已经注册"))
		return
	}
	if !this.dao.ChainRegisterRecordNotFound(param.TargetChainId, param.SourceChainId, targetContract.Address, 1) {
		this.ResponseError(c, CHAIN_REGISTER_EXISTS_ERROR, errors.New("链已经注册"))
		return
	}
	//注册一条链 source(1)->target(2)
	go this.registerOneChain(db, address, privateKey, source, param.TargetChainId, errChan, param.AnchorAddresses, param.SignConfirmCount, idString, pledge, 1)
	//注册另一条链 target->source
	go this.registerOneChain(db, address, privateKey, target, param.SourceChainId, errChan, param.AnchorAddresses, param.SignConfirmCount, idString, pledge, 2)
	errMsg := ""
	for i := 0; i < 2; i++ {
		err := <-errChan
		if err != nil {
			errMsg += err.Error()
		}
	}
	if errMsg != "" {
		db.Rollback()
		this.ResponseError(c, CHAIN_REGISTER_ERROR, errors.New(errMsg))
		return
	}
	db.Commit()
	this.echoSuccess(c, "链注册成功")
}
func (this *Controller) registerOneChain(db *gorm.DB, from common.Address, privateKey *ecdsa.PrivateKey, api *blockchain.Api, targetChainId uint, errChan chan error, strAnchorAddresses []string, signConfirmCount uint, anchorIds string, pledge *big.Int, role uint) {
	callerConfig := &blockchain.CallerConfig{
		From:       from,
		PrivateKey: privateKey,
		NetworkId:  api.GetNetworkId(),
	}
	chain, err := this.dao.GetChain(targetChainId)
	if err != nil {
		errChan <- err
		return
	}
	contract, err := this.dao.GetContractByChainId(api.GetChainId())
	if err != nil {
		errChan <- err
		return
	}
	anchorAddresses := make([]common.Address, 0, len(strAnchorAddresses))

	for _, v := range strAnchorAddresses {
		anchorAddresses = append(anchorAddresses, common.HexToAddress(v))
	}
	registerConfig := &blockchain.RegisterChainConfig{
		AbiData:          []byte(contract.Abi),
		ContractAddress:  common.HexToAddress(contract.Address),
		TargetNetworkId:  chain.NetworkId,
		AnchorAddresses:  anchorAddresses,
		SignConfirmCount: uint8(signConfirmCount),
		MaxValue:         pledge,
	}
	hash, err := api.RegisterChain(registerConfig, callerConfig)
	if err != nil {
		fmt.Println("RegisterChain ", err)
		errChan <- err
		return
	}
	register := &dao.ChainRegister{
		SourceChainId:   api.GetChainId(),
		TargetChainId:   targetChainId,
		TxHash:          hash,
		Confirm:         signConfirmCount,
		AnchorAddresses: anchorIds,
		Address:         contract.Address,
	}
	registerId, err := this.dao.CreateChainRegisterByTx(db, register)
	if err != nil {
		fmt.Println("CreateChainRegisterByTx ", err)
		errChan <- err
		return
	}
	go func(id uint, hash string, targetChainId uint) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 30
		i := 0
		for {
			<-ticker.C
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateChainRegisterStatus(id, int(receipt.Status))
				if err != nil {
					fmt.Println(err)
					continue
				}
				if int(receipt.Status) == 1 {
					chainRegister, err := this.dao.GetChainRegisterById(id)
					if err != nil {
						fmt.Println(err)
						continue
					}
					idStrings := strings.Split(chainRegister.AnchorAddresses, ",")
					for _, idStr := range idStrings {
						id, err := strconv.ParseUint(idStr, 10, 64)
						if err != nil {
							fmt.Println(err)
							continue
						}
						if role == 1 {
							//1 source
							this.dao.UpdateSourceStatus(uint(id), uint(receipt.Status))
						} else if role == 2 {
							//2 target
							this.dao.UpdateTargetStatus(uint(id), uint(receipt.Status))
						}
					}
				}
				break
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(registerId, hash, targetChainId)
	errChan <- nil
}

type ChainRegisterResult struct {
	TotalCount  int                      `json:"total_count"`  //总记录数
	CurrentPage int                      `json:"current_page"` //当前页数
	PageSize    int                      `json:"page_size"`    //页的大小
	PageData    []*dao.ChainRegisterView `json:"page_data"`    //页的数据
}

// @Summary 跨链管理(获取注册日志类表)
// @Tags ListChainRegister
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页，默认1"
// @Param page_size query string true "页的记录数，默认10"
// @Success 200 {object} JsonResult{data=ChainRegisterResult}
// @Router /chain/register/list [get]
func (this *Controller) ListChainRegister(c *gin.Context) {
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
	start := (currentPage - 1) * pageSize

	objects, err := this.dao.GetChainRegisterPage(start, pageSize)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	count, err := this.dao.GetChainRegisterCount()
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	chainRegisterResult := &ChainRegisterResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, chainRegisterResult)
}

type ChainRegisterInfo struct {
	dao.ChainRegisterView
	AnchorNodes []*dao.AnchorNode `json:"anchor_nodes" gorm:"anchor_nodes"`
}

// @Summary 获取注册日志详情
// @Tags GetChainRegisterInfo
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id query string true "注册日志id"
// @Success 200 {object} JsonResult{data=ChainRegisterInfo}
// @Router /chain/register/info [get]
func (this *Controller) GetChainRegisterInfo(c *gin.Context) {
	chainRegisterIdStr := c.Query("id")
	if chainRegisterIdStr == "" {
		this.ResponseError(c, REQUEST_PARAM_ERROR, errors.New("缺少参数 id"))
		return
	}
	id, err := strconv.ParseUint(chainRegisterIdStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	chain, err := this.dao.GetChainRegister(uint(id))
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	chainRegisterInfo := &ChainRegisterInfo{
		ChainRegisterView: *chain,
		AnchorNodes:       make([]*dao.AnchorNode, 0),
	}
	if chain.AnchorAddresses != "" {
		idStrings := strings.Split(chain.AnchorAddresses, ",")
		for _, idStr := range idStrings {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			anchorNode, err := this.dao.GetAnchorNode(uint(id))
			if err != nil {
				fmt.Println(err)
				continue
			}
			chainRegisterInfo.AnchorNodes = append(chainRegisterInfo.AnchorNodes, anchorNode)
		}
	}
	this.echoResult(c, chainRegisterInfo)
}
