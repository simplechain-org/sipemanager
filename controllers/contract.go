package controllers

import (
	"fmt"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
	"github.com/gin-gonic/gin"
	"sipemanager/blockchain"
	"sipemanager/dao"
	"strconv"
	"strings"
	"time"
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
	privateKey, err := blockchain.GetPrivateKey(wallet.Content, param.Password)
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
	privateKey, err := blockchain.GetPrivateKey(wallet.Content, param.Password)
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
