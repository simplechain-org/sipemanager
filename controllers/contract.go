package controllers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"

	"sipemanager/blockchain"
	"sipemanager/dao"
)

const (
	CONTRACT_INVOKE_ERROR  int = 12001 //合约调用出错
	CONTRACT_IN_USED_ERROR int = 12002 //合约正在使用
	CONTRACT_ID_NOT_EXISTS_ERROR int = 12003 //合约id对应的合约不存在
	CONTRACT_DEPLOY_ERROR int = 12004 //部署合约出错
	CONTRACT_CHECK_ERROR   int = 12005 //跨链合约Abi检查出错
)

// @Summary 上传本地合约
// @Tags AddContract
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param name formData string true "合约名称"
// @Param sol formData string true "合约源码"
// @Param abi formData string true "合约abi"
// @Param bin formData string true "合约bin"
// @Success 200 {object} JsonResult{data=int}
// @Router /contract/add [post]
func (this *Controller) AddContract(c *gin.Context) {
	var param dao.Contract
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	id, err := this.dao.CreateContract(&param)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, id)
}

type UpdateContractParam struct {
	Id   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Sol  string `gorm:"type:text" json:"sol" binding:"required"`
	Abi  string `gorm:"type:text" json:"abi" binding:"required"`
	Bin  string `gorm:"type:text" json:"bin" binding:"required"`
}

// @Summary 更新合约内容
// @Tags AddContract
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id formData string true "合约id"
// @Param name formData string true "合约名称"
// @Param sol formData string true "合约源码"
// @Param abi formData string true "合约abi"
// @Param bin formData string true "合约bin"
// @Success 200 {object} JsonResult{data=int}
// @Router /contract/update [put]
func (this *Controller) updateContract(c *gin.Context) {
	var param UpdateContractParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR, err)
		return
	}
	err := this.dao.UpdateContract(param.Id, param.Name, param.Sol, param.Abi, param.Bin)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR, err)
		return
	}
	this.echoSuccess(c, "Success")
}

// @Summary 删除合约
// @Tags RemoveContract
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path string true "合约id"
// @Success 200 {object} JsonResult{}
// @Router /contract/remove/{contract_id} [delete]
func (this *Controller) RemoveContract(c *gin.Context) {
	contractIdStr := c.Param("contract_id")
	if contractIdStr == "" {
		this.ResponseError(c, REQUEST_PARAM_ERROR, errors.New("缺少参数 contract_id"))
		return
	}
	contractId, err := strconv.ParseUint(contractIdStr, 10, 64)
	if err != nil {
		this.ResponseError(c, REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	can, err := this.dao.ContractCanDelete(uint(contractId))
	if err == nil && !can {
		this.ResponseError(c, CONTRACT_IN_USED_ERROR, errors.New("合约还在使用中，不可以删除"))
		return
	}
	err = this.dao.RemoveContract(uint(contractId))
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	this.echoSuccess(c, "Success")
}

type ContractResult struct {
	TotalCount  int             `json:"total_count"`  //总记录数
	CurrentPage int             `json:"current_page"` //当前页数
	PageSize    int             `json:"page_size"`    //页的大小
	PageData    []*dao.Contract `json:"page_data"`    //页的数据
}

// @Summary 合约管理
// @Tags ListContract
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页，默认1"
// @Param page_size query string true "页的记录数，默认10"
// @Success 200 {object} JsonResult{data=ContractResult}
// @Router /contract/list [get]
func (this *Controller) ListContract(c *gin.Context) {
	//分页
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
	status := c.Query("status")

	start := (currentPage - 1) * pageSize

	objects, err := this.dao.GetContractPage(start, pageSize, status)
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	count, err := this.dao.GetContractCount(status)
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	contractResult := &ContractResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, contractResult)

}

// @Summary 获取部署在链上的合约实例
// @Tags GetContractOnChain
// @Accept  json
// @Produce  json
// @Param chain_id query int true "链的id"
// @Security ApiKeyAuth
// @Success 200 {object} JsonResult{data=[]dao.ContractInstance}
// @Router /contract/chain [get]
func (this *Controller) GetContractOnChain(c *gin.Context) {
	chainId := c.Query("chain_id")
	if chainId == "" {
		this.ResponseError(c,REQUEST_PARAM_INVALID_ERROR, errors.New("chain_id非法"))
		return
	}
	id, err := strconv.ParseUint(chainId, 10, 64)
	if err != nil {
		this.ResponseError(c,REQUEST_PARAM_INVALID_ERROR, err)
		return
	}
	contracts, err := this.dao.GetContractsByChainId(uint(id))
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, contracts)
}

type ContractParam struct {
	Password   string `json:"password" form:"password" binding:"required"`
	WalletId   uint   `json:"wallet_id" form:"wallet_id" binding:"required"`
	ContractId uint   `json:"contract_id" form:"contract_id" binding:"required"`
	NodeId     uint   `json:"node_id" form:"node_id" binding:"required"`
}

// @Summary 本地合约上链
// @Tags InstanceContract
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param node_id formData string true "节点id"
// @Param contract_id formData string true "合约id"
// @Param wallet_id formData string true "钱包id"
// @Param password formData string true "钱包密码"
// @Success 200 {object} JsonResult{data=int}
// @Router /contract/instance [post]
func (this *Controller) InstanceContract(c *gin.Context) {
	var param ContractParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c,REQUEST_PARAM_ERROR, err)
		return
	}
	api, err := this.getApiByNodeId(param.NodeId)
	if err != nil {
		this.ResponseError(c,NODE_ID_EXISTS_ERROR, err)
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		this.ResponseError(c,WALLET_ID_NOT_EXISTS_ERROR, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		this.ResponseError(c,WALLET_PASSWORD_ERROR, err)
		return
	}
	contract, err := this.dao.GetContractById(param.ContractId)
	if err != nil {
		this.ResponseError(c,CONTRACT_ID_NOT_EXISTS_ERROR, err)
		return
	}
	var data []byte
	data = common.Hex2Bytes(contract.Bin)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	hash, err := api.DeployContract(address, nil, data, api.GetNetworkId(), privateKey)
	if err != nil {
		this.ResponseError(c, CONTRACT_DEPLOY_ERROR,err)
		return
	}
	contractObj := &dao.ContractInstance{
		TxHash:     hash,
		ChainId:    api.GetChainId(),
		ContractId: param.ContractId,
	}
	contractId, err := this.dao.CreateContractInstance(contractObj)
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	go func(id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 30
		i := 0
		for {
			<-ticker.C
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateContractAddress(contractId, receipt.ContractAddress.String())
				if err != nil {
					fmt.Println(err)
					continue
				}
				break
			}
			if i >= maxCount {
				break
			}
			i++
		}
	}(contractId, hash)

	this.echoResult(c, hash)
}

type ExistsContractParam struct {
	ChainId uint   `json:"chain_id"` //链id ,合约部署在那条链上
	TxHash  string `json:"tx_hash"`
	Address string `json:"address"`
	Name    string `json:"name" binding:"required"`
	Sol     string `json:"sol"`
	Abi     string `json:"abi" binding:"required"`
	Bin     string `json:"bin"`
}

//@Summary 引用链上合约
//@Accept application/x-www-form-urlencoded
//@Accept application/json
//@Produce application/json
//@Param chain_id formData uint true "链id"
//@Param tx_hash formData string true "交易哈希"
//@Param address formData string true "合约地址"
//@Param name formData string true "合约名称"
//@Param sol formData string true "合约源码"
//@Param abi formData string true "合约abi"
//@Param bin formData string bin "合约bin"
//@Success 200 {object} JsonResult
//@Router /contract/instance/import [post]
func (this *Controller) AddExistsContract(c *gin.Context) {
	var param ExistsContractParam
	if err := c.ShouldBind(&param); err != nil {
		this.ResponseError(c, REQUEST_PARAM_ERROR,err)
		return
	}
	contract := &dao.Contract{
		Name: param.Name,
		Sol:  param.Sol,
		Bin:  param.Bin,
		Abi:  param.Abi,
	}
	//创建合约对象
	id, err := this.dao.CreateContract(contract)
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	instance := &dao.ContractInstance{
		ChainId:    param.ChainId,
		TxHash:     param.TxHash,
		Address:    param.Address,
		ContractId: id,
	}
	id, err = this.dao.CreateContractInstance(instance)
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, "Success")
}
//@Summary 获取所有的合约
//@Accept application/x-www-form-urlencoded
//@Accept application/json
//@Produce application/json
//@Param chain_id formData uint true "链id"
//@Param tx_hash formData string true "交易哈希"
//@Param address formData string true "合约地址"
//@Param name formData string true "合约名称"
//@Param sol formData string true "合约源码"
//@Param abi formData string true "合约abi"
//@Param bin formData string bin "合约bin"
//@Success 200 {object} JsonResult
//@Router /contract/instance/import [post]
func (this *Controller) ListContractAll(c *gin.Context) {
	contracts, err := this.dao.GetContracts()
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	this.echoResult(c, contracts)
}

type ContractInstanceView struct {
	ID         uint   `json:"id"`
	CreatedAt  string `json:"created_at"`
	ChainId    uint   `gorm:"chain_id" json:"chain_id"` //链id ,合约部署在那条链上
	TxHash     string `gorm:"tx_hash" json:"tx_hash"`
	Address    string `gorm:"address" json:"address"`
	ContractId uint   `gorm:"contract_id" json:"contract_id"` //合约id
	Name       string `gorm:"name" json:"name"`
	ChainName  string `json:"chain_name"`
}
type ContractInstanceResult struct {
	TotalCount  int                    `json:"total_count"`  //总记录数
	CurrentPage int                    `json:"current_page"` //当前页数
	PageSize    int                    `json:"page_size"`    //页的大小
	PageData    []*dao.ContractInstanceView `json:"page_data"`    //页的数据
}

// @Summary 合约上链
// @Tags ListChain
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param current_page query string true "当前页，默认1"
// @Param page_size query string true "页的记录数，默认10"
// @Success 200 {object} JsonResult{data=ContractInstanceResult}
// @Router /contract/instance/list [get]
func (this *Controller) ListContractInstances(c *gin.Context) {
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

	objects, err := this.dao.GetContractInstancePage(start, pageSize)
	if err != nil {
		this.ResponseError(c, DATABASE_ERROR,err)
		return
	}
	count, err := this.dao.GetContractInstanceCount()
	if err != nil {
		this.ResponseError(c,DATABASE_ERROR, err)
		return
	}
	chainResult := &ContractInstanceResult{
		TotalCount:  count,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		PageData:    objects,
	}
	this.echoResult(c, chainResult)
}

const (
	MakerMethod    string = "2c50336e"
	MakerFinMethod string = "870f1f4a"
	TakerMethod    string = "f7478f6a"
)

func (this *Controller) CheckContractAbi(contractInstanceId uint) error {
	instance, err := this.dao.GetContractInstanceById(contractInstanceId)
	contract, err := this.dao.GetContractById(instance.ContractId)
	abiParsed, err := abi.JSON(strings.NewReader(contract.Abi))
	makerStart := hex.EncodeToString(abiParsed.Methods["makerStart"].ID())
	makerFinish := hex.EncodeToString(abiParsed.Methods["makerFinish"].ID())
	taker := hex.EncodeToString(abiParsed.Methods["taker"].ID())
	if makerStart != MakerMethod || makerFinish != MakerFinMethod || taker != TakerMethod {
		return errors.New("Unable to parse ABi normally")
	}
	if err != nil {
		return err
	}
	return nil
}
