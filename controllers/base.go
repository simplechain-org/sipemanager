package controllers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"sipemanager/blockchain"

	"github.com/gin-gonic/gin"
)

const (
	REQUEST_PARAM_ERROR int = 19001
	DATABASE_ERROR      int = 19300 //数据库错误
)

var dateFormat string = "2006-01-02 15:04:05"

type JsonResult struct {
	Msg  string      `json:"err_msg"`
	Code int         `json:"code" binding:"required"`
	Data interface{} `json:"data"`
}

func (this *Controller) echoError(c *gin.Context, err error) {
	logrus.Error(err)
	c.JSON(http.StatusOK, gin.H{
		"msg":  err.Error(),
		"code": -1,
		"data": nil,
	})
}
func (this *Controller) echoResult(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  "Success",
		"code": 0,
		"data": result,
	})
}
func (this *Controller) echoSuccess(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": 0,
		"data": nil,
	})
}

func (this *Controller) CheckHealth(c *gin.Context) {
	this.echoSuccess(c, "server is running")
}
func (this *Controller) getApiByNodeId(id uint) (*blockchain.Api, error) {
	node, err := this.dao.GetNodeById(id)
	if err != nil {
		return nil, err
	}
	chain, err := this.dao.GetChain(node.ChainId)
	if err != nil {
		return nil, err
	}
	n := &blockchain.Node{
		Address:   node.Address,
		Port:      node.Port,
		ChainId:   node.ChainId,
		IsHttps:   node.IsHttps,
		NetworkId: chain.NetworkId,
	}
	api, err := blockchain.NewApi(n)
	if err != nil {
		return nil, err
	}
	return api, nil
}
func (this *Controller) ResponseError(c *gin.Context, code int, err error) {
	logrus.Error(err)
	c.JSON(http.StatusOK, gin.H{
		"msg":  err.Error(),
		"code": code,
		"data": nil,
	})
}
