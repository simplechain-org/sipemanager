package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"sipemanager/dao"
)

type WalletParam struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type AddressParam struct {
	Address string `json:"address" binding:"required"`
}

// @Summary 添加钱包
// @Tags wallet
// @Accept  json
// @Produce  json
// @Param name formData string true "钱包昵称"
// @Param content formData string true "keystore string"
// @Security ApiKeyAuth
// @Success 200 {object} JSONResult{data=int} "walletId"
// @Router /wallet [post]
func (this *Controller) AddWallet(c *gin.Context) {
	var params WalletParam
	if err := c.ShouldBindJSON(&params); err != nil {
		this.echoError(c, err)
		return
	}
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var addressParam AddressParam
	err = json.Unmarshal([]byte(params.Content), &addressParam)
	if err != nil {
		this.echoError(c, err)
		return
	}
	address := "0x" + addressParam.Address
	wallet := dao.Wallet{
		Name:    params.Name,
		Content: []byte(params.Content),
		UserId:  user.ID,
		Address: address,
	}
	wallet.UserId = user.ID
	id, err := this.dao.CreateWallet(&wallet)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, id)
}

//不加载content
// @Summary 钱包列表
// @Tags wallet
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} JSONResult{data=dao.Wallet}
// @Router /wallet/list [get]
func (this *Controller) ListWallet(c *gin.Context) {
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	wallets, err := this.dao.ListWalletByUserId(user.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, wallets)
}
