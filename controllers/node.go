package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"sipemanager/dao"
)

//添加node
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

func (this *Controller) ChangeNode(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var params UserNodeParam
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
