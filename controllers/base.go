package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var dateFormat string = "2006-01-02 15:04:05"

type JsonResult struct {
	Msg  string      `json:"err_msg"`
	Code int         `json:"code" binding:"required"`
	Data interface{} `json:"data"`
}

func (this *Controller) echoError(c *gin.Context, err error) {
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
