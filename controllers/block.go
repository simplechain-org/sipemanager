package controllers

import (
	"math/big"
	"sipemanager/utils"
	"strconv"
	"strings"
	"sync"

	"sipemanager/blockchain"
	"sipemanager/dao"

	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	userClient  map[uint]*blockchain.Api
	dao         *dao.DataBaseAccessObject
	NodeChannel chan BlockChannel
	group       sync.WaitGroup
}

type Block struct {
	ParentHash   common.Hash    `json:"parentHash"`
	UncleHash    common.Hash    `json:"sha3Uncles"`
	CoinBase     common.Address `json:"miner"`
	Difficulty   *big.Int       `json:"difficulty"`
	Number       *big.Int       `json:"number"`
	GasLimit     uint64         `json:"gasLimit"`
	GasUsed      uint64         `json:"gasUsed"`
	Time         uint64         `json:"timestamp"`
	Nonce        uint64         `json:"nonce"`
	Transactions int            `json:"transactions"`
}
func (this *Controller) SyncBlock(api *blockchain.Api, number int64, node dao.InstanceNodes) {
	chainId := node.ChainId
	//fmt.Printf("----当前写入区块号:%+v, ----当前ChainId： %+v ————\n", number, chainId)
	block, err := api.BlockByNumber(big.NewInt(0).SetInt64(number))
	if err != nil {
		defer utils.DeferRecoverLog("controller => block => SyncBlock:", err.Error(), 30004, nil)
		panic(err.Error())
	}
	blockRecord := dao.Block{
		ParentHash:   block.ParentHash().Hex(),
		UncleHash:    block.UncleHash().Hex(),
		CoinBase:     block.Coinbase().Hex(),
		Difficulty:   block.Difficulty().Int64(),
		Number:       block.Number().Int64(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		Time:         block.Time(),
		Nonce:        block.Nonce(),
		Transactions: len(block.Transactions()),
		BlockHash:    block.Hash().Hex(),
		ChainId:      chainId,
	}
	replaceErr := this.dao.BlockReplace(blockRecord)
	if replaceErr != nil {
		logrus.Warn(utils.ErrLogCode{LogType: "controller => block => SyncBlock:", Code: 30003, Message: err.Error(), Err: nil})
	}

	if len(block.Transactions()) > 0 {
		for index, transaction := range block.Transactions() {
			netId := api.GetNetworkId()
			// 获取from地址
			msg, err := transaction.AsMessage(types.NewEIP155Signer(big.NewInt(0).SetUint64(netId)))
			if err != nil {
				logrus.Error(err.Error())
			}
			from := msg.From().Hex()
			receipt, err := api.TransactionReceipt(transaction.Hash())
			if err != nil {
				logrus.Error("TransactionReceipt:", err)
			}
			// hexutil.Encode(transaction.Data()),
			txRecord := dao.Transaction{
				BlockHash:        block.Hash().Hex(),
				BlockNumber:      block.Number().Int64(),
				Hash:             transaction.Hash().Hex(),
				From:             strings.ToLower(from),
				GasUsed:          strconv.FormatUint(block.GasUsed(), 10),
				GasPrice:         transaction.GasPrice().String(),
				Input:            transaction.Data(),
				Nonce:            transaction.Nonce(),
				TransactionIndex: int64(index),
				Value:            transaction.Value().String(),
				Timestamp:        block.Time(),
				Status:           receipt.Status,
				ChainId:          chainId,
			}
			if transaction.To() != nil {
				txRecord.To = strings.ToLower(transaction.To().Hex())
			}
			txReplaceErr := this.dao.TxReplace(txRecord)
			if txReplaceErr != nil {
				logrus.Error("Transactions Create:", txReplaceErr.Error())
			}
			contract, err := this.dao.GetContractById(node.ContractId)
			abiParsed, err := abi.JSON(strings.NewReader(contract.Abi))

			if len(transaction.Data()) < 4 {
				defer utils.DeferRecoverLog("dao => block => SyncBlock: ", "data len to low", 30002, nil)
				panic("data len to low")
			}
			sigdata := transaction.Data()[:4]
			argdata := transaction.Data()[4:]
			method, err := abiParsed.MethodById(sigdata)
			if err != nil {
				defer utils.DeferRecoverLog("dao => block => SyncBlock: ", err.Error(), 30003, err)
				panic(err.Error())
			}
			var args dao.MakerFinish
			UnpackErr := method.Inputs.Unpack(&args, argdata)
			if UnpackErr != nil {

			} else {
				txRecord.EventType = "makerFinish"
				txReplaceErr := this.dao.TxReplace(txRecord)
				if txReplaceErr != nil {
					logrus.Error("Transactions Create:", txReplaceErr.Error())
				}
				targetChain, err := this.dao.GetChainByNetWorkId(args.RemoteChainId.Uint64())
				if err != nil {
					logrus.Error("CrossAnchors => GetChainByNetWorkId:", err.Error())
				}
				crossRecord := dao.CrossAnchors{
					BlockNumber:     txRecord.BlockNumber,
					GasUsed:         txRecord.GasUsed,
					GasPrice:        txRecord.GasPrice,
					ContractAddress: txRecord.To,
					Timestamp:       txRecord.Timestamp,
					Status:          txRecord.Status,
					ChainId:         txRecord.ChainId,
					RemoteChainId:   targetChain.ID,
					EventType:       txRecord.EventType,
					NetworkId:       node.NetworkId,
					RemoteNetworkId: args.RemoteChainId.Uint64(),
					AnchorAddress:   txRecord.From,
					TxId:            args.Rtx.TxId.Hex(),
					Hash:            txRecord.Hash,
				}
				crossErr := this.dao.CrossAnchorsReplace(crossRecord)
				if crossErr != nil {
					logrus.Error("CrossAnchors Create:", crossErr.Error())
				}
			}
		}
	}

	if len(block.Uncles()) > 0 {
		for _, uncle := range block.Uncles() {
			uncleRecord := dao.Uncle{
				ParentHash:  uncle.ParentHash.Hex(),
				UncleHash:   uncle.UncleHash.Hex(),
				CoinBase:    uncle.Coinbase.Hex(),
				Difficulty:  uncle.Difficulty.Int64(),
				Number:      uncle.Number.Int64(),
				GasLimit:    uncle.GasLimit,
				GasUsed:     uncle.GasUsed,
				Time:        uncle.Time,
				Nonce:       uncle.Nonce.Uint64(),
				ChainId:     chainId,
				BlockNumber: block.Number().Int64(),
			}
			uncleReplaceErr := this.dao.UncleReplace(uncleRecord)
			if uncleReplaceErr != nil {
				logrus.Error("UncleRecord Create:", err.Error())
			}
		}

	}
}
func (this *Controller) getApi(userId uint, networkId uint64) (*blockchain.Api, error) {
	node, err := this.dao.GetNodeByUserIdAndNetworkId(userId, networkId)
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
