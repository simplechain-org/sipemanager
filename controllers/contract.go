package controllers

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"

	"sipemanager/blockchain"
	"sipemanager/dao"
)

type Contract struct {
	Password string `json:"password" binding:"required"`
	WalletId uint   `json:"wallet_id" binding:"required"`
	Sol      string `json:"sol"`
	Abi      string `json:"abi" binding:"required"`
	Bin      string `json:"bin" binding:"required"`
}

func (this *Controller) ListContract(c *gin.Context) {
	contracts, err := this.dao.GetContracts()
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, contracts)
}
func (this *Controller) GetContractOnChain(c *gin.Context) {
	chainId := c.Query("chain_id")
	if chainId != "" {
		id, err := strconv.ParseUint(chainId, 10, 64)
		if err != nil {
			this.echoError(c, err)
			return
		}
		contracts, err := this.dao.GetContractsByChainId(uint(id))
		if err != nil {
			this.echoError(c, err)
			return
		}
		this.echoResult(c, contracts)
		return
	}
}
func (this *Controller) GetContractInstances(c *gin.Context) {
	contracts, err := this.dao.GetContractInstances()
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, contracts)
	return
}

type ContractParam struct {
	Password   string `json:"password" binding:"required"`
	WalletId   uint   `json:"wallet_id" binding:"required"`
	ContractId uint   `json:"contract_id" binding:"required"`
}

func (this *Controller) AddContract(c *gin.Context) {
	var param dao.Contract
	if err := c.ShouldBindJSON(&param); err != nil {
		this.echoError(c, err)
		return
	}
	id, err := this.dao.CreateContract(&param)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, id)
}

//部署合约
func (this *Controller) DeployContract(c *gin.Context) {
	var param ContractParam
	if err := c.ShouldBindJSON(&param); err != nil {
		this.echoError(c, err)
		return
	}
	user, err := this.GetUser(c)
	if err != nil {
		this.echoError(c, err)
		return
	}
	api, err := this.getBlockChainApi(user.ID)
	if err != nil {
		this.echoError(c, err)
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		this.echoError(c, err)
		return
	}
	contract, err := this.dao.GetContractById(param.ContractId)
	var data []byte
	data = common.Hex2Bytes(contract.Bin)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	hash, err := api.DeployContract(address, nil, data, api.GetNetworkId(), privateKey)
	if err != nil {
		this.echoError(c, err)
		return
	}
	contractObj := &dao.ContractInstance{
		TxHash:     hash,
		ChainId:    api.GetChainId(),
		ContractId: param.ContractId,
	}
	contractId, err := this.dao.CreateContractInstance(contractObj)
	if err != nil {
		this.echoError(c, err)
		return
	}

	go func(id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
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

type RegisterChainParam struct {
	TargetChainId    uint     `json:"target_chain_id"`
	SignConfirmCount uint     `json:"sign_confirm_count"`
	AnchorAddresses  []string `json:"anchor_addresses"`
	WalletId         uint     `json:"wallet_id"`
	Password         string   `json:"password"`
}

func (this *Controller) ListRegisterChain(c *gin.Context) {
	result, err := this.dao.ListChainRegister()
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, result)
}

func (this *Controller) RegisterChain(c *gin.Context) {
	var param RegisterChainParam
	if err := c.ShouldBindJSON(&param); err != nil {
		fmt.Println("ShouldBindJSON:", err.Error())
		this.echoError(c, err)
		return
	}
	user, err := this.GetUser(c)
	if err != nil {
		fmt.Println("GetUser:", err.Error())
		this.echoError(c, err)
		return
	}
	api, err := this.getBlockChainApi(user.ID)
	if err != nil {
		fmt.Println("getBlockChainApi:", err.Error())
		this.echoError(c, err)
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		fmt.Println("GetWallet:", err.Error())
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		fmt.Println("GetPrivateKey:", err.Error())
		this.echoError(c, err)
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	callerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  api.GetNetworkId(),
	}
	chain, err := this.dao.GetChain(param.TargetChainId)
	if err != nil {
		fmt.Println("GetChain:", err.Error())
		this.echoError(c, err)
		return
	}
	anchorAddresses := make([]common.Address, 0, len(param.AnchorAddresses))

	for _, v := range param.AnchorAddresses {
		anchorAddresses = append(anchorAddresses, common.HexToAddress(v))
	}
	contract, err := this.dao.GetContractByChainId(api.GetChainId())

	if err != nil {
		fmt.Println("GetContractByChainId:", err.Error())
		this.echoError(c, err)
		return
	}

	registerConfig := &blockchain.RegisterChainConfig{
		AbiData:          []byte(contract.Abi),
		ContractAddress:  common.HexToAddress(contract.Address),
		TargetNetworkId:  uint64(chain.NetworkId),
		AnchorAddresses:  anchorAddresses,
		SignConfirmCount: uint8(param.SignConfirmCount),
	}
	hash, err := api.RegisterChain(registerConfig, callerConfig)
	if err != nil {
		fmt.Println("RegisterChain:", err.Error())
		this.echoError(c, err)
		return
	}
	register := &dao.ChainRegister{
		SourceChainId:   api.GetChainId(),
		TargetChainId:   param.TargetChainId,
		TxHash:          hash,
		Confirm:         param.SignConfirmCount,
		AnchorAddresses: strings.Join(param.AnchorAddresses, ","),
	}
	registerId, err := this.dao.CreateChainRegister(register)
	if err != nil {
		fmt.Println("CreateChainRegister:", err.Error())
		this.echoError(c, err)
		return
	}
	go func(id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateChainRegisterStatus(id, int(receipt.Status))
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
	}(registerId, hash)

	this.echoResult(c, hash)
}

//todo 判重
func (this *Controller) AddContractInstance(c *gin.Context) {
	var param dao.ContractInstance
	if err := c.ShouldBindJSON(&param); err != nil {
		this.echoError(c, err)
		return
	}
	contractId, err := this.dao.CreateContractInstance(&param)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, contractId)
}

type RegisterChainAddParam struct {
	SourceChainId    uint     `json:"source_chain_id"`
	TargetChainId    uint     `json:"target_chain_id"`
	SignConfirmCount uint     `json:"sign_confirm_count"`
	AnchorAddresses  []string `json:"anchor_addresses"`
	TxHash           string   `json:"tx_hash"`
}

func (this *Controller) RegisterChainAdd(c *gin.Context) {
	var param RegisterChainAddParam
	if err := c.ShouldBindJSON(&param); err != nil {
		this.echoError(c, err)
		return
	}
	register := &dao.ChainRegister{
		SourceChainId:   param.SourceChainId,
		TargetChainId:   param.TargetChainId,
		TxHash:          param.TxHash,
		Confirm:         param.SignConfirmCount,
		AnchorAddresses: strings.Join(param.AnchorAddresses, ","),
		Status:          1,
		StatusText:      "成功",
	}
	registerId, err := this.dao.CreateChainRegister(register)
	if err != nil {
		fmt.Println("CreateChainRegister:", err.Error())
		this.echoError(c, err)
		return
	}
	this.echoResult(c, registerId)
}

type ExistsContractParam struct {
	ChainId     uint   `json:"chain_id"` //链id ,合约部署在那条链上
	TxHash      string `json:"tx_hash"`
	Address     string `json:"address"`
	Description string `json:"description" binding:"required"`
	Sol         string `json:"sol"`
	Abi         string `json:"abi" binding:"required"`
	Bin         string `json:"bin"`
}

//引用链上合约

//@Summary 引用链上合约
//@Description 引用链上合约
//@Accept application/x-www-form-urlencoded
//@Accept application/json
//@Produce application/json
//@Param "" body ExistsContractParam true "请求体"
//@Success 200 {object} JsonResult
//@Router /api/v1/contract/instance/import [post]
func (this *Controller) AddExistsContract(c *gin.Context) {
	var param ExistsContractParam
	if err := c.ShouldBind(&param); err != nil {
		this.echoError(c, err)
		return
	}
	contract := &dao.Contract{
		Description: param.Description,
		Sol:         param.Sol,
		Bin:         param.Bin,
		Abi:         param.Abi,
	}
	//创建合约对象
	id, err := this.dao.CreateContract(contract)
	if err != nil {
		this.echoError(c, err)
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
		this.echoError(c, err)
		return
	}
	this.echoResult(c, "success")
}

type RegisterChainTwoWayParam struct {
	SourceChainId    uint     `json:"source_chain_id" form:"source_chain_id"`
	TargetChainId    uint     `json:"target_chain_id" form:"target_chain_id"`
	SourceNodeId     uint     `json:"source_node_id" form:"source_node_id"`
	TargetNodeId     uint     `json:"target_node_id" form:"target_node_id"`
	SignConfirmCount uint     `json:"sign_confirm_count" form:"sign_confirm_count"`
	AnchorAddresses  []string `json:"anchor_addresses" form:"anchor_addresses"`
	AnchorNames      []string `json:"anchor_names" form:"anchor_names"`
	WalletId         uint     `json:"wallet_id" form:"wallet_id"`
	Password         string   `json:"password" form:"password"`
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

//@Summary 引用链上合约
//@Description 引用链上合约
//@Accept application/x-www-form-urlencoded
//@Accept application/json
//@Produce application/json
//@Param "" body ExistsContractParam true "请求体"
//@Success 200 {object} JsonResult
//@Router /api/v1/contract/register/once [post]
func (this *Controller) RegisterChainTwoWay(c *gin.Context) {
	var param RegisterChainTwoWayParam
	if err := c.ShouldBind(&param); err != nil {
		fmt.Println("ShouldBind:", err.Error())
		this.echoError(c, err)
		return
	}
	wallet, err := this.dao.GetWallet(param.WalletId)
	if err != nil {
		fmt.Println("GetWallet:", err.Error())
		this.echoError(c, err)
		return
	}
	privateKey, err := blockchain.GetPrivateKey([]byte(wallet.Content), param.Password)
	if err != nil {
		fmt.Println("GetPrivateKey:", err.Error())
		this.echoError(c, err)
		return
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	errChan := make(chan error, 2)

	source, err := this.getApiByNodeId(param.SourceNodeId)
	if err != nil {
		fmt.Println("GetApiByNodeId:", err.Error())
		this.echoError(c, err)
		return
	}
	target, err := this.getApiByNodeId(param.TargetNodeId)
	if err != nil {
		fmt.Println("GetApiByNodeId:", err.Error())
		this.echoError(c, err)
		return
	}
	db := this.dao.BeginTransaction()
	for index, address := range param.AnchorAddresses {
		anchorNode := &dao.AnchorNode{
			Name:          param.AnchorNames[index],
			Address:       address,
			SourceChainId: param.SourceChainId,
			TargetChainId: param.TargetChainId,
		}
		_, err := this.dao.CreateAnchorNodeByTx(db, anchorNode)
		if err != nil {
			db.Rollback()
			this.echoError(c, err)
			return
		}
	}
	//注册一条链 source->target
	go this.registerOneChain(db, address, privateKey, source, param.TargetChainId, errChan, param.AnchorAddresses, param.SignConfirmCount)
	//注册另一条链 target->source
	go this.registerOneChain(db, address, privateKey, target, param.SourceChainId, errChan, param.AnchorAddresses, param.SignConfirmCount)

	errMsg := ""
	for i := 0; i < 2; i++ {
		err := <-errChan
		if err != nil {
			errMsg += err.Error()
		}
	}
	if errMsg != "" {
		db.Rollback()
		this.echoError(c, errors.New(errMsg))
		return
	}
	db.Commit()
	this.echoSuccess(c, "链注册成功")
}

func (this *Controller) registerOneChain(db *gorm.DB, from common.Address, privateKey *ecdsa.PrivateKey, api *blockchain.Api, targetChainId uint, errChan chan error, strAnchorAddresses []string, signConfirmCount uint) {
	callerConfig := &blockchain.CallerConfig{
		From:       from,
		PrivateKey: privateKey,
		NetworkId:  api.GetNetworkId(),
	}
	chain, err := this.dao.GetChain(targetChainId)
	if err != nil {
		errChan <- err
		return
	}
	contract, err := this.dao.GetContractByChainId(api.GetChainId())
	if err != nil {
		errChan <- err
		return
	}
	anchorAddresses := make([]common.Address, 0, len(strAnchorAddresses))

	for _, v := range strAnchorAddresses {
		anchorAddresses = append(anchorAddresses, common.HexToAddress(v))
	}
	registerConfig := &blockchain.RegisterChainConfig{
		AbiData:          []byte(contract.Abi),
		ContractAddress:  common.HexToAddress(contract.Address),
		TargetNetworkId:  uint64(chain.NetworkId),
		AnchorAddresses:  anchorAddresses,
		SignConfirmCount: uint8(signConfirmCount),
	}
	hash, err := api.RegisterChain(registerConfig, callerConfig)
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println("hash=", hash)
	register := &dao.ChainRegister{
		SourceChainId:   api.GetChainId(),
		TargetChainId:   targetChainId,
		TxHash:          hash,
		Confirm:         signConfirmCount,
		AnchorAddresses: strings.Join(strAnchorAddresses, ","),
		Address:         contract.Address,
	}
	registerId, err := this.dao.CreateChainRegisterByTx(db, register)
	if err != nil {
		errChan <- err
		return
	}
	go func(id uint, hash string) {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		maxCount := 300 //最多尝试300次
		i := 0
		for {
			<-ticker.C
			fmt.Println("now:", time.Now().Unix())
			//时间到，做一下检测
			receipt, err := api.TransactionReceipt(common.HexToHash(hash))
			if err == nil && receipt != nil {
				err = this.dao.UpdateChainRegisterStatus(id, int(receipt.Status))
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
	}(registerId, hash)
}
