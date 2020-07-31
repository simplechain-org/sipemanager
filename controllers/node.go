package controllers

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"sipemanager/dao"

	"github.com/gin-gonic/gin"
)

const (
	NODE_ID_EXISTS_ERROR int = 11001 //节点id对应的记录不存在
)

//添加node
// @Summary 添加node
// @Tags AddNode
// @Accept  json
// @Produce  json
// @Param chain_id formData int true "关联链Id"
// @Param address formData string true "节点地址"
// @Param port formData int true "端口"
// @Param name formData string true "名称"
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=int} "NodeId"
// @Router /node [post]
func (this *Controller) AddNode(c *gin.Context) {
	var node dao.Node
	if err := c.ShouldBind(&node); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	if node.Port < 0 || node.Port > 65535 {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("port参数非法"))
		return
	}
	trial := net.ParseIP(node.Address)
	if trial == nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("不是一个ip地址"))
		return
	}
	if trial.To4() == nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, fmt.Errorf("%v is not an IPv4 address\n", trial))
		return
	}
	node.UserId = user.ID
	chain, err := this.dao.GetChain(node.ChainId)
	if err != nil {
		this.ResponseError(c, CHAIN_ID_NOT_EXISTS_ERROR, err)
		return
	}
	node.ChainId = chain.ID
	id, err := this.dao.CreateNode(&node)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, id)
}

type UpdateNodeParam struct {
	Id      uint   `json:"id" binding:"required" form:"id"`
	Address string `gorm:"size:255" json:"address" binding:"required" form:"address"` //地址
	Port    int    `json:"port" binding:"required" form:"port"`                       //端口
	Name    string `gorm:"size:255" json:"name" binding:"required" form:"name"`
	ChainId uint   `json:"chain_id" binding:"required" form:"chain_id"` //链id
}

// @Summary 编辑节点
// @Tags UpdateNode
// @Accept  json
// @Produce  json
// @Param id formData int true "Id"
// @Param chain_id formData int true "关联链Id"
// @Param address formData string true "节点地址"
// @Param port formData int true "端口"
// @Param name formData string true "名称"
// @Security ApiKeyAuth
// @Success 200 {object}
// @Router /node [put]
func (this *Controller) UpdateNode(c *gin.Context) {
	var params UpdateNodeParam
	if err := c.ShouldBind(&params); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	if params.Port < 0 || params.Port > 65535 {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("port参数非法"))
		return
	}
	trial := net.ParseIP(params.Address)
	if trial == nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("不是一个ip地址"))
		return
	}
	if trial.To4() == nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, fmt.Errorf("%v is not an IPv4 address\n", trial))
		return
	}
	if params.Name == "" {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("name参数非法"))
		return
	}
	err := this.dao.UpdateNode(params.Id, params.Address, params.Port, params.Name, params.ChainId)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	this.echoSuccess(c, "Success")
}

// @Summary 删除节点
// @Tags DeleteNode
// @Accept  json
// @Produce  json
// @Param id path int true "Id"
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{msg=string}
// @Router /node/remove/{id} [delete]
func (this *Controller) DeleteNode(c *gin.Context) {
	nodeIdStr := c.Param("id")
	if nodeIdStr == "" {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, errors.New("缺少参数 id"))
		return
	}
	nodeId, err := strconv.ParseUint(nodeIdStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	err = this.dao.RemoveNode(uint(nodeId))
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoSuccess(c, "Success")
}

// @Summary 节点列表(获取用户的所有节点)
// @Tags GetAllNodes
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=[]dao.NodeView}
// @Router /node/list/all [get]
func (this *Controller) GetAllNodes(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	nodes, err := this.dao.ListNodeByUserId(user.ID)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, nodes)
}

type NodeResult struct {
	TotalCount  int            `json:"total_count"`  //总记录数
	CurrentPage int            `json:"current_page"` //当前页数
	PageSize    int            `json:"page_size"`    //页的大小
	PageData    []dao.NodeView `json:"page_data"`    //页的数据
}

// @Summary 节点的管理（分页获取节点列表）
// @Tags ListNode
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页，默认1"
// @Param page_size query string true "页的记录数，默认10"
// @Success 200 {object} JsonResult{data=NodeResult}
// @Router /node/list/page [get]
func (this *Controller) ListNode(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
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
	objects, err := this.dao.ListNodeByUserIdPage(start, pageSize, user.ID)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	count, err := this.dao.GetNodeByUserIdCount(user.ID)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	chainResult := &NodeResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, chainResult)
}
