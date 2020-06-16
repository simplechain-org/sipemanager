package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (this *Controller) echoError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  err.Error(),
		"code": -1,
	})
}
func (this *Controller) echoResult(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"result": result,
	})
}

func (this *Controller) CheckHealth(c *gin.Context) {
	this.echoResult(c, "server is running")
}
