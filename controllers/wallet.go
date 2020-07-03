package controllers

import (
	"encoding/json"

	"sipemanager/dao"
	"sipemanager/utils"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

type WalletParam struct {
	Name     string `json:"name" binding:"required" form:"name"`         //账户名称
	Content  string `json:"content" binding:"required" form:"content"`   //私钥/助记词/keystore文件
	Password string `json:"password" binding:"required" form:"password"` //钱包密码
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
	var err error
	if err := c.ShouldBindJSON(&params); err != nil {
		this.echoError(c, err)
		return
	}
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var address string
	var content []byte

	if utils.IsHex(params.Content) {
		//私钥
		privateKeyECDSA, err := crypto.HexToECDSA(params.Content)
		if err != nil {
			this.echoError(c, err)
			return
		}
		content, err = utils.PrivateKeyToKeystore(privateKeyECDSA, params.Password)
		if err != nil {
			this.echoError(c, err)
			return
		}
		address, err = utils.GetAddress(privateKeyECDSA)
		if err != nil {
			this.echoError(c, err)
			return
		}
	} else {
		if utils.IsJSON(params.Content) {
			//keystore文件内容
			//校验口令
			_, err := keystore.DecryptKey([]byte(params.Content), params.Password)
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
			address = "0x" + addressParam.Address
			content = []byte(params.Content)
		} else {
			//助记词
			privateKeyECDSA, err := utils.GetPrivateKeyFromMnemonic(params.Content)
			if err != nil {
				this.echoError(c, err)
				return
			}
			content, err = utils.PrivateKeyToKeystore(privateKeyECDSA, params.Password)
			if err != nil {
				this.echoError(c, err)
				return
			}
			address, err = utils.GetAddress(privateKeyECDSA)
			if err != nil {
				this.echoError(c, err)
				return
			}
		}
	}
	wallet := dao.Wallet{
		Name:    params.Name,
		Content: content,
		UserId:  user.ID,
		Address: address,
	}
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

type UpdateWalletParam struct {
	WalletId    uint   `json:"wallet_id" binding:"required" form:"wallet_id"`       //钱包id
	OldPassword string `json:"old_password" binding:"required" form:"old_password"` //钱包旧密码
	NewPassword string `json:"new_password" binding:"required" form:"new_password"` //钱包新密码
}

func (this *Controller) UpdateWallet(c *gin.Context) {
	var params UpdateWalletParam
	if err := c.ShouldBind(&params); err != nil {
		this.echoError(c, err)
		return
	}
	wallet, err := this.dao.GetWallet(params.WalletId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	key, err := keystore.DecryptKey([]byte(wallet.Content), params.OldPassword)
	if err != nil {
		this.echoError(c, err)
		return
	}
	content, err := utils.PrivateKeyToKeystore(key.PrivateKey, params.NewPassword)
	if err != nil {
		this.echoError(c, err)
		return
	}
	err = this.dao.UpdateWallet(params.WalletId, content)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoSuccess(c, "success")
}
