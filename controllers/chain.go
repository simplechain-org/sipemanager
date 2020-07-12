package controllers

import (
	"errors"
	"strconv"

	"sipemanager/dao"

	"github.com/gin-gonic/gin"
)

type ChainContractParam struct {
	ChainId    uint `json:"chain_id" binding:"required"`
	ContractId uint `json:"contract_instance_id" binding:"required"`
}

//获取所有的链信息
func (this *Controller) GetChains(c *gin.Context) {
	chains, err := this.dao.GetChains()
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, chains)
}

//因为是异步得到地址的，所以需要独立更新地址
//更新链的跨链合约地址
func (this *Controller) UpdateChainContractAddress(c *gin.Context) {
	var param ChainContractParam
	if err := c.ShouldBindJSON(&param); err != nil {
		this.echoError(c, err)
		return
	}
	err := this.dao.UpdateChainContract(&dao.ChainContract{
		ChainId:            param.ChainId,
		ContractInstanceId: param.ContractId,
	})
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, "success")
}

//获取链信息
func (this *Controller) GetChainInfo(c *gin.Context) {
	chainIdStr := c.Param("chain_id")
	if chainIdStr == "" {
		this.echoError(c, errors.New("缺少参数 chain_id"))
		return
	}
	chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
	if err != nil {
		this.echoError(c, err)
		return
	}
	chain, err := this.dao.GetChain(uint(chainId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, chain)
}

//保存链信息
func (this *Controller) CreateChain(c *gin.Context) {
	var param dao.Chain
	if err := c.ShouldBindJSON(&param); err != nil {
		this.echoError(c, err)
		return
	}
	id, err := this.dao.CreateChain(&param)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, id)
}

//删除链信息
// @Summary 删除链信息
// @Tags chain
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param chain_id path int true "chainId"
// @Success 200 {object} JsonResult{data=string} "成功删除链信息"
// @Router /chain/{chain_id} [delete]
func (this *Controller) RemoveChain(c *gin.Context) {
	chainIdStr := c.Param("chain_id")
	if chainIdStr == "" {
		this.echoError(c, errors.New("缺少参数 chain_id"))
		return
	}
	chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
	if err != nil {
		this.echoError(c, err)
		return
	}
	check := this.dao.ChainHasNode(uint(chainId))
	if check {
		this.echoError(c, errors.New("还存在相应的节点，不能删除"))
		return
	}
	check = this.dao.ChainHasContractInstance(uint(chainId))
	if check {
		this.echoError(c, errors.New("还存在相应的合约实例，不能删除"))
		return
	}
	err = this.dao.ChainRemove(uint(chainId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, "成功删除链信息")
}
