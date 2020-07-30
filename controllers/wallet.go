package controllers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"

	"sipemanager/dao"
	"sipemanager/utils"

	"github.com/gin-gonic/gin"
	"github.com/simplechain-org/go-simplechain/accounts/keystore"
	"github.com/simplechain-org/go-simplechain/crypto"
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
// @Param password formData string true "密码"
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=int} "walletId"
// @Router /wallet [post]
func (this *Controller) AddWallet(c *gin.Context) {
	var params WalletParam
	var err error
	if err := c.ShouldBind(&params); err != nil {
		this.echoError(c, err)
		return
	}
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	var address string
	var content string

	if utils.IsHex(params.Content) {
		//私钥
		privateKeyECDSA, err := crypto.HexToECDSA(params.Content)
		if err != nil {
			this.echoError(c, err)
			return
		}
		keyData, err := utils.PrivateKeyToKeystore(privateKeyECDSA, params.Password)
		if err != nil {
			this.echoError(c, err)
			return
		}
		address, err = utils.GetAddress(privateKeyECDSA)
		if err != nil {
			this.echoError(c, err)
			return
		}
		content = string(keyData)
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
			content = string(params.Content)
		} else {
			//助记词
			privateKeyECDSA, err := utils.GetPrivateKeyFromMnemonic(params.Content)
			if err != nil {
				this.echoError(c, err)
				return
			}
			KeyData, err := utils.PrivateKeyToKeystore(privateKeyECDSA, params.Password)
			if err != nil {
				this.echoError(c, err)
				return
			}
			address, err = utils.GetAddress(privateKeyECDSA)
			if err != nil {
				this.echoError(c, err)
				return
			}
			content = string(KeyData)
		}
	}
	if this.dao.WalletExists(address) {
		this.echoError(c, errors.New("钱包地址已经存在"))
		return
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

// @Summary 钱包列表
// @Tags wallet
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=dao.Wallet}
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

//修改钱包密码（口令）
// @Summary 修改钱包密码
// @Tags updateWallet
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param wallet_id formData string true "钱包id"
// @Param old_password formData string true "钱包旧密码"
// @Param new_password formData string true "钱包新密码"
// @Success 200 {object} JsonResult{data=object}
// @Router /wallet/update [post]
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
	this.echoSuccess(c, "Success")
}

// @Summary 移除钱包
// @Tags removeWallet
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param wallet_id query string true "钱包id"
// @Success 200 {object} JsonResult{msg=string}
// @Router /wallet/remove [delete]
func (this *Controller) RemoveWallet(c *gin.Context) {
	walletIdStr := c.Query("wallet_id")
	walletId, err := strconv.ParseUint(walletIdStr, 10, 64)
	if err != nil {
		this.echoError(c, err)
		return
	}
	_, err = this.dao.GetWallet(uint(walletId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	err = this.dao.RemoveWallet(uint(walletId))
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoSuccess(c, "Success")
}
type WalletResult struct {
	TotalCount  int              `json:"total_count"`  //总记录数
	CurrentPage int              `json:"current_page"` //当前页数
	PageSize    int              `json:"page_size"`    //页的大小
	PageData    []*dao.WalletView `json:"page_data"`    //页的数据
}
// @Summary 钱包列表(分页显示)
// @Tags wallet
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页"
// @Param page_size query string true "页的记录数"
// @Success 200 {object} JsonResult{data=WalletResult}
// @Router /wallet/list/page [get]
func (this *Controller) ListPageWallet(c *gin.Context) {
	var pageSize int = 10
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
	//当前页（默认为第一页）
	var currentPage int = 1
	currentPageStr := c.Query("current_page")
	if currentPageStr != "" {
		page, err := strconv.ParseUint(currentPageStr, 10, 64)
		if err == nil {
			currentPage = int(page)
		}
	}
	start := (currentPage - 1) * pageSize
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	wallets, err := this.dao.GetWalletViewPage(user.ID, start, pageSize)
	if err != nil {
		this.echoError(c, err)
		return
	}
	count, err := this.dao.GetWalletViewCount(user.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	walletResult := &WalletResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageData:    wallets,
		PageSize:    pageSize,
	}
	this.echoResult(c, walletResult)
}


