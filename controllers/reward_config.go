package controllers

import (
	"errors"
	"math/big"
	"strconv"

	"sipemanager/dao"

	"github.com/gin-gonic/gin"
)

// @Summary 配置签名奖励
// @Tags AddRewardConfig
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param source_chain_id formData uint true "发起链id"
// @Param target_chain_id formData uint true "目标链id"
// @Param regulation_cycle formData uint true "调控周期"
// @Param sign_reward formData string true "单笔签名奖励"
// @Success 200 {object} JsonResult{data=int}
// @Router /reward/config/add [post]
func (this *Controller) AddRewardConfig(c *gin.Context) {
	var param dao.RewardConfig
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	id, err := this.dao.CreateRewardConfig(&param)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, id)
}

// @Summary 获取配置签名奖励详情
// @Tags GetRewardConfigInfo
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path uint true "签名奖励id"
// @Success 200 {object} JsonResult{data=dao.RewardConfigView}
// @Router /reward/config/info/:id [get]
func (this *Controller) GetRewardConfigInfo(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		this.ResponseError(c, REQUEST_PARAM_ERROR, errors.New("缺少参数 id"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	result, err := this.dao.GetRewardConfig(uint(id))
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, result)
}

// @Summary 删除签名奖励
// @Tags AddRewardConfig
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "签名奖励id"
// @Success 200 {object}  JsonResult{}
// @Router /reward/config/remove/{id} [delete]
func (this *Controller) RemoveRewardConfig(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		this.ResponseError(c, REQUEST_PARAM_ERROR, errors.New("缺少参数 id"))
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	err = this.dao.RemoveRelativeRewardConfig(uint(id))
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoSuccess(c, "Success")
}

type RewardConfigResult struct {
	TotalCount  int                     `json:"total_count"`  //总记录数
	CurrentPage int                     `json:"current_page"` //当前页数
	PageSize    int                     `json:"page_size"`    //页的大小
	PageData    []*dao.RewardConfigView `json:"page_data"`    //页的数据
}

// @Summary 配置签名奖励 列表
// @Tags ListRewardConfig
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页，默认1"
// @Param page_size query string true "页的记录数，默认10"
// @Success 200 {object} JsonResult{data=RewardConfigResult}
// @Router /reward/config/list [get]
func (this *Controller) ListRewardConfig(c *gin.Context) {
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

	objects, err := this.dao.GetRewardConfigPage(start, pageSize)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	count, err := this.dao.GetRewardConfigCount()
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	rewardConfigResult := &RewardConfigResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, rewardConfigResult)
}

type GetRewardConfigParam struct {
	SourceChainId uint `gorm:"source_chain_id" json:"source_chain_id"`
	TargetChainId uint `gorm:"target_chain_id" json:"target_chain_id"`
}

// @Summary 获取配置签名奖励详情
// @Tags GetRewardConfig
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param source_chain_id formData uint true "发起链id"
// @Param target_chain_id formData uint true "目标链id"
// @Success 200 {object} JsonResult{data=dao.RewardConfigView}
// @Router /reward/config/detail [post]
func (this *Controller) GetRewardConfig(c *gin.Context) {
	var param GetRewardConfigParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	result, err := this.dao.GetLatestRewardConfig(param.SourceChainId, param.TargetChainId)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, result)
}

type UpdateRewardConfigParam struct {
	Id              uint   `json:"id"`
	RegulationCycle uint   `json:"regulation_cycle"`
	SignReward      string `json:"sign_reward"`
}

// @Summary 更新签名奖励
// @Tags UpdateRewardConfig
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id formData uint true "签名奖励id"
// @Param regulation_cycle formData uint true "调控周期"
// @Param sign_reward formData string true "单笔签名奖励"
// @Success 200 {object} JsonResult{data=object}
// @Router /reward/config/update [post]
func (this *Controller) UpdateRewardConfig(c *gin.Context) {
	var param UpdateRewardConfigParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	last, err := this.dao.GetRewardConfigById(param.Id)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	//调整周期必须大于0
	if param.RegulationCycle == 0 {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("调控周期必须大于0"))
		return
	}
	//校验sign_reward为整数
	_, ok := big.NewInt(0).SetString(param.SignReward, 10)
	if !ok {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("sign_reward数据非法"))
		return
	}
	rewardConfig := &dao.RewardConfig{
		SourceChainId:   last.SourceChainId,
		TargetChainId:   last.TargetChainId,
		RegulationCycle: param.RegulationCycle,
		SignReward:      param.SignReward,
	}
	//记录历史，所以每次更新都是新增一条记录
	id, err := this.dao.CreateRewardConfig(rewardConfig)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, id)
}
