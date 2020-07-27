package controllers

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/gin-gonic/gin"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
)

type AddPrepareRewardParam struct {
	SourceChainId uint   `json:"source_chain_id" binding:"required"`
	TargetChainId uint   `json:"target_chain_id" binding:"required"`
	SourceReward  string `json:"source_reward" binding:"required"`
	TargetReward  string `json:"target_reward" binding:"required"`
	WalletId      uint   `json:"wallet_id" form:"wallet_id"`
	Password      string `json:"password" form:"password"`
}

//@Summary 添加跨链手续费
//@Accept application/x-www-form-urlencoded
//@Accept application/json
//@Produce application/json
//@Param source_chain_id formData uint true "源链Id"
//@Param target_chain_id formData uint true "目标链Id"
//@Param source_reward formData string true "源链跨链手续费"
//@Param target_reward formData string true "目标链跨链手续费"
//@Param wallet_id formData uint true "钱包id"
//@Param password formData string true "钱包密码"
//@Success 200 {object} JsonResult{data=int}
//@Router /chain/cross/prepare/reward [post]
func (this *Controller) AddPrepareReward(c *gin.Context) {
	var param AddPrepareRewardParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, errors.New("参数错误:"+err.Error()))
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		this.echoError(c, errors.New("钱包不存在"))
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		this.echoError(c, errors.New("密码错误"))
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	sourceChain, err := this.dao.GetChain(param.SourceChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("不存在链 id=%d", param.SourceChainId))
		return
	}
	sourceNode, err := this.dao.GetNodeByChainId(param.SourceChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("chain=%d下不存在节点", param.SourceChainId))
		return
	}
	sourceContract, err := this.dao.GetContractByChainId(param.SourceChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("chain=%d下不存在跨链合约", param.SourceChainId))
		return
	}
	targetChain, err := this.dao.GetChain(param.TargetChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("不存在链 id=%d", param.TargetChainId))
		return
	}
	targetNode, err := this.dao.GetNodeByChainId(param.TargetChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("chain=%d下不存在节点", param.TargetChainId))
		return
	}
	targetContract, err := this.dao.GetContractByChainId(param.TargetChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("chain=%d下不存在跨链合约", param.TargetChainId))
		return
	}
	sourceReward, success := big.NewInt(0).SetString(param.SourceReward, 10)
	if !success {
		this.echoError(c, errors.New("source_reward数据非法"))
		return
	}
	targetReward, success := big.NewInt(0).SetString(param.TargetReward, 10)
	if !success {
		this.echoError(c, errors.New("target_reward数据非法"))
		return
	}
	//链的合约
	sourceConfig := &blockchain.SetRewardConfig{
		AbiData:         []byte(sourceContract.Abi),
		ContractAddress: common.HexToAddress(sourceContract.Address),
		TargetNetworkId: targetChain.NetworkId,
	}
	sourceCallerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  sourceChain.NetworkId,
	}
	sourceApi, err := this.getApiByNodeId(sourceNode.ID)
	if err != nil {
		this.echoError(c, errors.New("为节点构建api失败"))
		return
	}
	sourceHash, err := sourceApi.SetReward(sourceConfig, sourceCallerConfig, sourceReward)
	if err != nil {
		this.echoError(c, errors.New("请求设置链上数据失败:"+err.Error()))
		return
	}
	targetConfig := &blockchain.SetRewardConfig{
		AbiData:         []byte(targetContract.Abi),
		ContractAddress: common.HexToAddress(targetContract.Address),
		TargetNetworkId: sourceChain.NetworkId,
	}
	targetCallerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  targetChain.NetworkId,
	}
	targetApi, err := this.getApiByNodeId(targetNode.ID)
	if err != nil {
		this.echoError(c, errors.New("构建api失败:"+err.Error()))
		return
	}
	targetHash, err := targetApi.SetReward(targetConfig, targetCallerConfig, targetReward)
	if err != nil {
		this.echoError(c, errors.New("设置链上数据失败:"+err.Error()))
		return
	}
	prepareReward := &dao.PrepareReward{
		SourceChainId: param.SourceChainId,
		TargetChainId: param.TargetChainId,
		SourceReward:  param.SourceReward,
		TargetReward:  param.TargetReward,
		SourceHash:    sourceHash,
		TargetHash:    targetHash,
	}
	id, err := this.dao.CreatePrepareReward(prepareReward)
	if err != nil {
		this.echoError(c, errors.New("保存预扣手续费出错"))
		return
	}
	this.echoResult(c, id)
}

//@Summary 更新跨链手续费
//@Accept application/x-www-form-urlencoded
//@Accept application/json
//@Produce application/json
//@Param source_chain_id formData uint true "源链Id"
//@Param target_chain_id formData uint true "目标链Id"
//@Param source_reward formData string true "源链跨链手续费"
//@Param target_reward formData string true "目标链跨链手续费"
//@Param wallet_id formData uint true "钱包id"
//@Param password formData string true "钱包密码"
//@Success 200 {object} JsonResult
//@Router /chain/cross/prepare/reward [put]
type UpdatePrepareRewardParam struct {
	SourceChainId uint   `json:"source_chain_id" binding:"required"`
	TargetChainId uint   `json:"target_chain_id" binding:"required"`
	SourceReward  string `json:"source_reward" binding:"required"`
	TargetReward  string `json:"target_reward" binding:"required"`
	WalletId      uint   `json:"wallet_id" form:"wallet_id"`
	Password      string `json:"password" form:"password"`
}

func (this *Controller) UpdatePrepareReward(c *gin.Context) {
	var param UpdatePrepareRewardParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, errors.New("参数错误:"+err.Error()))
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		this.echoError(c, errors.New("钱包不存在"))
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		this.echoError(c, errors.New("密码错误"))
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	sourceChain, err := this.dao.GetChain(param.SourceChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("不存在链 id=%d", param.SourceChainId))
		return
	}
	sourceNode, err := this.dao.GetNodeByChainId(param.SourceChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("chain=%d下不存在节点", param.SourceChainId))
		return
	}
	sourceContract, err := this.dao.GetContractByChainId(param.SourceChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("chain=%d下不存在跨链合约", param.SourceChainId))
		return
	}
	targetChain, err := this.dao.GetChain(param.TargetChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("不存在链 id=%d", param.TargetChainId))
		return
	}
	targetNode, err := this.dao.GetNodeByChainId(param.TargetChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("chain=%d下不存在节点", param.TargetChainId))
		return
	}
	targetContract, err := this.dao.GetContractByChainId(param.TargetChainId)
	if err != nil {
		this.echoError(c, fmt.Errorf("chain=%d下不存在跨链合约", param.TargetChainId))
		return
	}
	sourceReward, success := big.NewInt(0).SetString(param.SourceReward, 10)
	if !success {
		this.echoError(c, errors.New("source_reward数据非法"))
		return
	}
	targetReward, success := big.NewInt(0).SetString(param.TargetReward, 10)
	if !success {
		this.echoError(c, errors.New("target_reward数据非法"))
		return
	}
	//链的合约
	sourceConfig := &blockchain.SetRewardConfig{
		AbiData:         []byte(sourceContract.Abi),
		ContractAddress: common.HexToAddress(sourceContract.Address),
		TargetNetworkId: targetChain.NetworkId,
	}
	sourceCallerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  sourceChain.NetworkId,
	}
	sourceApi, err := this.getApiByNodeId(sourceNode.ID)
	if err != nil {
		this.echoError(c, errors.New("为节点构建api失败"))
		return
	}
	sourceHash, err := sourceApi.SetReward(sourceConfig, sourceCallerConfig, sourceReward)
	if err != nil {
		this.echoError(c, errors.New("请求设置链上数据失败:"+err.Error()))
		return
	}
	targetConfig := &blockchain.SetRewardConfig{
		AbiData:         []byte(targetContract.Abi),
		ContractAddress: common.HexToAddress(targetContract.Address),
		TargetNetworkId: sourceChain.NetworkId,
	}
	targetCallerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  targetChain.NetworkId,
	}
	targetApi, err := this.getApiByNodeId(targetNode.ID)
	if err != nil {
		this.echoError(c, errors.New("构建api失败:"+err.Error()))
		return
	}
	targetHash, err := targetApi.SetReward(targetConfig, targetCallerConfig, targetReward)
	if err != nil {
		this.echoError(c, errors.New("设置链上数据失败:"+err.Error()))
		return
	}
	prepareReward := &dao.PrepareReward{
		SourceChainId: param.SourceChainId,
		TargetChainId: param.TargetChainId,
		SourceReward:  param.SourceReward,
		TargetReward:  param.TargetReward,
		SourceHash:    sourceHash,
		TargetHash:    targetHash,
	}
	err = this.dao.UpdatePrepareReward(prepareReward)
	if err != nil {
		this.echoError(c, errors.New("保存预扣手续费出错"))
		return
	}
	this.echoResult(c, "Success")
}
type PrepareRewardResult struct {
	TotalCount  int                     `json:"total_count"`  //总记录数
	CurrentPage int                     `json:"current_page"` //当前页数
	PageSize    int                     `json:"page_size"`    //页的大小
	PageData    []*dao.PrepareRewardView `json:"page_data"`    //页的数据
}
// @Summary 跨链配置（预扣手续费） 列表
// @Tags ListPrepareReward
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页，默认1"
// @Param page_size query string true "页的记录数，默认10"
// @Success 200 {object} JsonResult{data=PrepareRewardResult}
// @Router /reward/prepare/reward/list [get]
func (this *Controller) ListPrepareReward(c *gin.Context) {
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

	objects, err := this.dao.GetPrepareRewardPage(start,pageSize)
	if err != nil {
		this.echoError(c, err)
		return
	}
	count, err := this.dao.GetPrepareRewardCount()
	if err != nil {
		this.echoError(c, err)
		return
	}
	prepareRewardResult := &PrepareRewardResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, prepareRewardResult)
}
