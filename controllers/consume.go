package controllers

import (
	"fmt"
	"time"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/crypto"
	"github.com/gin-gonic/gin"

	"sipemanager/blockchain"
	"sipemanager/dao"
)

func (this *Controller) ListTakerOrder(c *gin.Context) {
	result, err := this.dao.ListTakerOrder()
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, result)
}

//罗列现存的跨链交易
func (this *Controller) ListCrossTransaction(c *gin.Context) {
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
	transactions, err := api.GetRemote(api.GetNetworkId())
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, transactions)
}

type ConsumeParam struct {
	CtxId    string `json:"ctx_id" binding:"required"`
	Password string `json:"password" binding:"required"`
	WalletId uint   `json:"wallet_id" binding:"required"`
}

//买入跨链交易（接单）
func (this *Controller) Consume(c *gin.Context) {
	var param ConsumeParam
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
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	callerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  api.GetNetworkId(),
	}
	//根据当前的链id获取到当前使用的合约
	contract, err := this.dao.GetContractByChainId(api.GetChainId())
	if err != nil {
		this.echoError(c, err)
		return
	}
	contractConfig := &blockchain.ContractConfig{
		AbiData:         []byte(contract.Abi),                  //合约的api
		ContractAddress: common.HexToAddress(contract.Address), //合约的地址
	}
	handleOrder, err := api.Consume(param.CtxId, contractConfig, callerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}
	targetChainId, err := this.dao.GetTargetChainId(api.GetChainId(), handleOrder.NetworkId)
	if err != nil {
		this.echoError(c, err)
		return
	}
	order := &dao.TakerOrder{
		SourceChainId: api.GetChainId(),
		TargetChainId: targetChainId,
		Taker:         address.String(),
		SourceValue:   handleOrder.SourceValue.Uint64(),
		TargetValue:   handleOrder.TargetValue.Uint64(),
		TxHash:        handleOrder.TxHash,
		CtxId:         handleOrder.CtxId,
	}

	orderId, err := this.dao.CreateTakerOrder(order)

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
				err = this.dao.UpdateTakerOrderStatus(id, int(receipt.Status))
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
	}(orderId, handleOrder.TxHash)

	this.echoResult(c, handleOrder.TxHash)
}
