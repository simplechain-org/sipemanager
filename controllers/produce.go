package controllers

import (
	"fmt"
	"math/big"
	"time"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

type ProduceParam struct {
	ChainId     uint   `json:"chain_id" binding:"required"`
	Password    string `json:"password" binding:"required"`
	WalletId    uint   `json:"wallet_id" binding:"required"`
	SourceValue uint64 `json:"source_value" binding:"required"`
	TargetValue uint64 `json:"target_value" binding:"required"`
	Extra       string `json:"extra"`
}

func (this *Controller) ListMakerOrder(c *gin.Context) {
	result, err := this.dao.ListMakerOrder()
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, result)
}

//测试通过
func (this *Controller) Produce(c *gin.Context) {
	var param ProduceParam
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
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	callerConfig := &blockchain.CallerConfig{
		From:       address,
		PrivateKey: privateKey,
		NetworkId:  api.GetNetworkId(),
	}
	contract, err := this.dao.GetContractByChainId(api.GetChainId())
	if err != nil {
		this.echoError(c, err)
		return
	}
	chain, err := this.dao.GetChain(param.ChainId)
	contractConfig := &blockchain.ContractConfig{
		AbiData:         []byte(contract.Abi),
		ContractAddress: common.HexToAddress(contract.Address),
		TargetChainId:   uint64(chain.NetworkId),
	}
	changeParam := &blockchain.ChangeParam{
		SourceValue: big.NewInt(0).SetUint64(param.SourceValue),
		TargetValue: big.NewInt(0).SetUint64(param.TargetValue),
		Input:       []byte(param.Extra),
	}
	hash, err := api.Produce(contractConfig, changeParam, callerConfig)
	if err != nil {
		this.echoError(c, err)
		return
	}

	//SourceChainId uint
	//TargetChainId uint
	//Maker        string
	//SourceValue   string //仅仅做记录用，不计算
	//TargetValue   string
	//Status        int
	//StatusText    string
	//TxHash        string

	maker := &dao.MakerOrder{
		SourceChainId: api.GetChainId(),
		TargetChainId: param.ChainId,
		Maker:         address.String(),
		SourceValue:   param.SourceValue,
		TargetValue:   param.TargetValue,
		TxHash:        hash,
	}
	makerId, err := this.dao.CreateMakerOrder(maker)

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
				err = this.dao.UpdateMakerOrderStatus(id, int(receipt.Status))
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
	}(makerId, hash)

	this.echoResult(c, hash)
}
