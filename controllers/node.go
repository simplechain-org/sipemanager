package controllers

import (
	"errors"
	"fmt"

	"sipemanager/dao"

	"github.com/gin-gonic/gin"
)

//添加node
// @Summary 添加node
// @Tags node
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
	if err := c.ShouldBindJSON(&node); err != nil {
		this.echoError(c, err)
		return
	}
	//todo 校验数据
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	node.UserId = user.ID
	chain, err := this.dao.GetChain(node.ChainId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	node.NetworkId = chain.NetworkId

	id, err := this.dao.CreateNode(&node)
	if err != nil {
		this.echoError(c, err)
		return
	}
	//第一次添加成功，我们设置它为默认节点
	if !this.dao.UserHasNode(user.ID) {
		_, err := this.dao.CreateUserNode(&dao.UserNode{
			UserId: user.ID,
			NodeId: id,
		})
		if err != nil {
			this.echoError(c, err)
			return
		}
	}
	this.echoResult(c, id)
}
func (this *Controller) UpdateNode(c *gin.Context) {

}
func (this *Controller) DeleteNode(c *gin.Context) {

}

//切换节点
type UserNodeParam struct {
	UserId uint `json:"user_id"`
	NodeId uint `json:"node_id"`
}

//切换node
// @Summary 切换node
// @Tags node
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param user_id formData int true "useId"
// @Param node_id formData int true "nodeId"
// @Success 200 {object} string "success"
// @Router /node/change [post]
func (this *Controller) ChangeNode(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var params UserNodeParam
	fmt.Printf("fdf%+v", params)
	if err := c.Bind(&params); err != nil {
		this.echoError(c, err)
		return
	}
	if user.ID != params.UserId {
		this.echoError(c, errors.New("invalid operation"))
		return
	}
	err = this.dao.UpdateUserCurrentNode(&dao.UserNode{UserId: user.ID, NodeId: params.NodeId})
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.onChangeNode(user.ID)
	this.echoResult(c, "success")
}

type Node struct {
	dao.Node
	Description string `json:"description"`
	ChainName   string `json:"chain_name"`
}

// @Summary 节点列表
// @Tags node
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=[]Node}
// @Router /node/list [get]
func (this *Controller) GetNodes(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	nodes, err := this.dao.ListNodeByUserId(user.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	result := make([]Node, 0)
	for _, node := range nodes {
		chain, err := this.dao.GetChain(node.ChainId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		description := fmt.Sprintf("%s_%s_%s:%d", chain.Name, node.Name, node.Address, node.Port)
		result = append(result, Node{Node: node, Description: description, ChainName: chain.Name})
	}
	this.echoResult(c, result)
}

// @Summary 获取当前登录账户的节点
// @Tags node
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=[]Node}
// @Router /node/current [get]
func (this *Controller) GetUserCurrentNode(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetUserCurrentNode(user.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	chain, err := this.dao.GetChain(node.ChainId)
	if err != nil {
		fmt.Println(err)
		return
	}
	description := fmt.Sprintf("%s_%s_%s:%d", chain.Name, node.Name, node.Address, node.Port)
	result := Node{Node: *node, Description: description}
	this.echoResult(c, result)
}
func (this *Controller) GetUserCurrentChain(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	node, err := this.dao.GetUserCurrentNode(user.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	nodeInfo, err := this.dao.GetNodeById(node.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, nodeInfo)
}
