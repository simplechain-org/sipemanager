package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"sipemanager/dao"

	"github.com/gin-gonic/gin"
)

const (
	CHAIN_SAVE_ERROR                int = 10001 //链信息保存出错
	CHAIN_UPDATE_ERROR              int = 10002 //更新链信息出错
	CHAIN_ID_NOT_EXISTS_ERROR       int = 10003 //链id对应的记录不存在
	CHAIN_CONTRACT_NOT_EXISTS_ERROR int = 10004 //链id对应的合约记录不存在
	CHAIN_HAD_NODE_ERROR            int = 10005 //链下还存在节点，不能删除
	CHAIN_HAD_CONTRACT_ERROR        int = 10006 //链下还存在相应的合约实例
	CHAIN_REGISTER_ERROR            int = 10007 //链注册出错
	CHAIN_REGISTER_EXISTS_ERROR     int = 10008 //链已经注册

)

// @Summary 添加链信息
// @Tags CreateChain
// @Accept  json
// @Produce  json
// @Param name formData int true "链的名称"
// @Param network_id formData string true "链的网络编号"
// @Param coin_name formData int true "币的名称"
// @Param symbol formData string true "币的符号"
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=int}
// @Router /chain/create [post]
func (this *Controller) CreateChain(c *gin.Context) {
	var param dao.Chain
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, fmt.Errorf("数据类型不匹配:%s", err.Error()))
		return
	}
	id, err := this.dao.CreateChain(&param)
	if err != nil {
		this.ResponseError(c, CHAIN_SAVE_ERROR, errors.New("保存链的基本信息时发生错误"))
		return
	}
	this.echoResult(c, id)
}

type UpdateChainParam struct {
	Id                 uint   `json:"id" binding:"required"`
	Name               string `json:"name" binding:"required"`       //链的名称
	NetworkId          uint64 `json:"network_id" binding:"required"` //链的网络编号
	CoinName           string `json:"coin_name" binding:"required"`  //币名
	Symbol             string `json:"symbol" binding:"required"`     //符号
	ContractInstanceId uint   `json:"contract_instance_id"`          //合约实例
}

// @Summary 编辑链信息
// @Tags node
// @Accept  json
// @Produce  json
// @Param id formData int true "链的id"
// @Param name formData int true "链的名称"
// @Param network_id formData string true "链的网络编号"
// @Param coin_name formData int true "币的名称"
// @Param symbol formData string true "币的符号"
// @Param contract_instance_id formData uint true "合约实例id"
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=int}
// @Router /chain/update [put]
func (this *Controller) UpdateChain(c *gin.Context) {
	var param UpdateChainParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, fmt.Errorf("数据类型不匹配:%s", err.Error()))
		return
	}
	if param.ContractInstanceId != 0 {
		checkErr := this.CheckContractAbi(param.ContractInstanceId)
		if checkErr != nil {
			this.ResponseError(c, CONTRACT_CHECK_ERROR, fmt.Errorf("跨链合约abi检查错误:%s", checkErr.Error()))
			return
		}
	}
	err := this.dao.UpdateChain(param.Id,
		param.Name, param.NetworkId, param.CoinName, param.Symbol, param.ContractInstanceId)
	if err != nil {
		this.ResponseError(c, CHAIN_UPDATE_ERROR, fmt.Errorf("更新链信息时发生错误:%s", err.Error()))
		return
	}
	if param.ContractInstanceId != 0 {
		go this.UpdateDirectBlock(param.Id)
	}
	this.echoSuccess(c, "Success")
}

// @Summary 删除链信息
// @Tags chain
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param chain_id path int true "chainId"
// @Success 200 {object} JsonResult{data=object} "成功删除链信息"
// @Router /chain/{chain_id} [delete]
func (this *Controller) RemoveChain(c *gin.Context) {
	chainIdStr := c.Param("chain_id")
	if chainIdStr == "" {
		this.ResponseError(c, REQUEST_PARAM_ERROR, errors.New("缺少参数 chain_id"))
		return
	}
	chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, fmt.Errorf("chain_id不是一个整数:%s", err.Error()))
		return
	}
	check := this.dao.ChainHasNode(uint(chainId))
	if check {
		this.ResponseError(c, CHAIN_HAD_NODE_ERROR, errors.New("还存在相应的节点，不能删除"))
		return
	}
	check = this.dao.ChainHasContractInstance(uint(chainId))
	if check {
		this.ResponseError(c, CHAIN_HAD_CONTRACT_ERROR, errors.New("还存在相应的合约实例，不能删除"))
		return
	}
	err = this.dao.ChainRemove(uint(chainId))
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, fmt.Errorf("删除链信息时发生错误:%s", err.Error()))
		return
	}
	this.echoResult(c, "成功删除链信息")
}

type ChainResult struct {
	TotalCount  int              `json:"total_count"`  //总记录数
	CurrentPage int              `json:"current_page"` //当前页数
	PageSize    int              `json:"page_size"`    //页的大小
	PageData    []*dao.ChainInfo `json:"page_data"`    //页的数据
}

// @Summary 链的管理
// @Tags ListChain
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页，默认1"
// @Param page_size query string true "页的记录数，默认10"
// @Success 200 {object} JsonResult{data=ChainResult}
// @Router /chain/list [get]
func (this *Controller) ListChain(c *gin.Context) {
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

	objects, err := this.dao.GetChainInfoPage(start, pageSize)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	count, err := this.dao.GetChainInfoCount()
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	chainResult := &ChainResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, chainResult)
}

// @Summary 获取链信息
// @Tags GetChainInfo
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param chain_id path string true "链id"
// @Success 200 {object} JsonResult{data=dao.ChainInfo}
// @Router /chain/info/{chain_id} [get]
func (this *Controller) GetChainInfo(c *gin.Context) {
	chainIdStr := c.Param("chain_id")
	if chainIdStr == "" {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("缺少参数 chain_id"))
		return
	}
	chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	chain, err := this.dao.GetChain(uint(chainId))
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, chain)
}

// @Summary 获取链相关的节点
// @Tags GetNodeByChain
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param chain_id query string true "链id"
// @Success 200 {object} JsonResult{data=[]dao.Node}
// @Router /chain/node [get]
func (this *Controller) GetNodeByChain(c *gin.Context) {
	chainIdStr := c.Query("chain_id")
	if chainIdStr == "" {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("缺少参数 chain_id"))
		return
	}
	chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	chain, err := this.dao.ListNodeByChainId(uint(chainId))
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, chain)
}

// @Summary 获取所有链信息
// @Tags ListAllChain
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=dao.Chain}
// @Router /chain/list/all [get]
func (this *Controller) ListAllChain(c *gin.Context) {
	chains, err := this.dao.ListAllChain()
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, chains)
}
