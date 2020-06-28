package controllers

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

type Transaction struct {
	AccountNonce uint64          `json:"nonce"`
	Price        *big.Int        `json:"gasPrice""`
	GasLimit     uint64          `json:"gas"`
	Recipient    *common.Address `json:"to"`
	Amount       *big.Int        `json:"value"`
	Payload      []byte          `json:"input"`
	Hash         string          `json:"hash"`
	From         common.Address  `json:"from"`
}

//分页列表
// @Summary 分页列表
// @Tags block
// @Accept  json
// @Produce  json
// @Param number query int true "blockNumber"
// @Success 200 {object} JSONResult{data=Transaction}
// @Router /block/transaction/{number} [post]
func (this *Controller) GetBlockTransaction(c *gin.Context) {
	numberStr := c.Param("number")
	number, err := strconv.ParseUint(numberStr, 10, 64)
	if err != nil {
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
	block, err := api.BlockByNumber(big.NewInt(0).SetUint64(number))
	if err != nil {
		fmt.Println(err)
		return
	}
	result := make([]Transaction, 0)

	signer := types.NewEIP155Signer(big.NewInt(int64(api.GetNetworkId())))

	for _, tx := range block.Transactions() {
		from, err := types.Sender(signer, tx)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, Transaction{
			AccountNonce: tx.Nonce(),
			Price:        tx.GasPrice(),
			GasLimit:     tx.Gas(),
			Recipient:    tx.To(),
			Amount:       tx.Value(),
			Hash:         tx.Hash().String(),
			From:         from,
		})

	}
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"result": result,
	})
}

type TransactionReceipt struct {
	TxHash            common.Hash     `json:"transactionHash"`
	Status            string          `json:"status"`
	AccountNonce      uint64          `json:"nonce"`
	Price             *big.Int        `json:"gasPrice"`
	GasLimit          uint64          `json:"gas"`
	Recipient         *common.Address `json:"to"`
	Amount            *big.Int        `json:"value"`
	PostState         []byte          `json:"root"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	ContractAddress   string          `json:"contractAddress"`
	GasUsed           uint64          `json:"gasUsed"`
	BlockHash         common.Hash     `json:"blockHash"`
	BlockNumber       *big.Int        `json:"blockNumber"`
	TransactionIndex  uint            `json:"transactionIndex"`
	Payload           string          `json:"input"`
}

func (this *Controller) GetTransactionReceipt(c *gin.Context) {
	hash := c.Param("hash")
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
	transaction, err := api.TransactionByHash(common.HexToHash(hash))
	if err != nil {
		this.echoError(c, err)
		return
	}
	receipt, err := api.TransactionReceipt(common.HexToHash(hash))
	if err != nil {
		this.echoError(c, err)
		return
	}
	var status string
	if receipt.Status == 1 {
		status = "success"
	} else {
		status = "failure"
	}
	var address string
	if (receipt.ContractAddress != common.Address{}) {
		address = receipt.ContractAddress.String()
	}
	result := TransactionReceipt{
		Status:            status,
		TxHash:            receipt.TxHash,
		AccountNonce:      transaction.Nonce(),
		Price:             transaction.GasPrice(),
		GasLimit:          transaction.Gas(),
		Recipient:         transaction.To(),
		Amount:            transaction.Value(),
		Payload:           common.Bytes2Hex(transaction.Data()),
		PostState:         receipt.PostState,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		ContractAddress:   address,
		GasUsed:           receipt.GasUsed,
		BlockHash:         receipt.BlockHash,
		BlockNumber:       receipt.BlockNumber,
		TransactionIndex:  receipt.TransactionIndex,
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"result": result,
	})
}
