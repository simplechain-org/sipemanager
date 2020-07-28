package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/simplechain-org/go-simplechain"
	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/common/hexutil"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/simplechain-org/go-simplechain/cross/backend"
	"github.com/simplechain-org/go-simplechain/rlp"
)

type CrossTransaction struct {
	Value            *big.Int       `json:"value"`
	CTxId            common.Hash    `json:"ctx_id"`
	TxHash           common.Hash    `json:"tx_hash"`
	From             common.Address `json:"from"`
	BlockHash        common.Hash    `json:"block_hash"`
	DestinationId    *big.Int       `json:"destination_id"`
	DestinationValue *big.Int       `json:"destination_value"`
	Input            hexutil.Bytes  `json:"input"`
}
type RPCCrossTransaction struct {
	Value            *hexutil.Big   `json:"Value"`
	CTxId            common.Hash    `json:"ctxId"`
	TxHash           common.Hash    `json:"TxHash"`
	From             common.Address `json:"From"`
	To               common.Address `json:"to"`
	BlockHash        common.Hash    `json:"BlockHash"`
	DestinationId    *hexutil.Big   `json:"destinationId"`
	DestinationValue *hexutil.Big   `json:"DestinationValue"`
	Input            hexutil.Bytes  `json:"input"`
	V                []*hexutil.Big `json:"V"`
	R                []*hexutil.Big `json:"R"`
	S                []*hexutil.Big `json:"S"`
}

func (this *Api) GetRemote(destinationId uint64) ([]*CrossTransaction, error) {
	var signatures map[string]map[uint64][]*RPCCrossTransaction
	err := this.client.CallContext(context.Background(), &signatures, "eth_ctxContent")
	if err != nil {
		return nil, err
	}
	result := make([]*CrossTransaction, 0)
	for _, value := range signatures["remote"] {
		for index, _ := range value {
			v := value[index]
			fmt.Println("v=", v.DestinationId.ToInt().Int64())
			if v.DestinationId.ToInt().Int64() == int64(destinationId) {
				result = append(result, &CrossTransaction{
					Value:            v.Value.ToInt(),
					CTxId:            v.CTxId,
					TxHash:           v.TxHash,
					From:             v.From,
					BlockHash:        v.BlockHash,
					DestinationId:    v.DestinationId.ToInt(),
					DestinationValue: v.DestinationValue.ToInt(),
				})
			}
		}
	}
	return result, nil
}

type RegisterChainConfig struct {
	AbiData          []byte
	ContractAddress  common.Address
	TargetNetworkId  uint64
	AnchorAddresses  []common.Address
	SignConfirmCount uint8
	MaxValue         *big.Int
}

type CallerConfig struct {
	From       common.Address
	PrivateKey *ecdsa.PrivateKey
	NetworkId  uint64
}

func (this *Api) RegisterChain(config *RegisterChainConfig, callerConfig *CallerConfig) (string, error) {
	abi, err := abi.JSON(bytes.NewReader(config.AbiData))
	if err != nil {
		return "", err
	}
	out, err := abi.Pack("chainRegister", big.NewInt(0).SetUint64(config.TargetNetworkId),config.MaxValue, config.SignConfirmCount, config.AnchorAddresses)

	if err != nil {
		return "", err
	}
	nonce, err := this.simpleClient.PendingNonceAt(context.Background(), callerConfig.From)
	if err != nil {
		return "", err
	}
	gasPrice, err := this.simpleClient.SuggestGasPrice(context.Background())

	if err != nil {
		return "", err
	}
	msg := simplechain.CallMsg{
		From:     callerConfig.From,
		To:       &config.ContractAddress,
		Data:     out,
		GasPrice: gasPrice,
	}
	gasLimit, err := this.simpleClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", err
	}

	transaction := types.NewTransaction(nonce, config.ContractAddress, big.NewInt(0), gasLimit, gasPrice, out)
	transaction, err = types.SignTx(transaction, types.NewEIP155Signer(big.NewInt(0).SetInt64(int64(callerConfig.NetworkId))), callerConfig.PrivateKey)
	if err != nil {
		return "", err
	}
	content, err := rlp.EncodeToBytes(transaction)

	if err != nil {
		return "", err
	}
	var result common.Hash

	err = this.client.CallContext(context.Background(), &result, "eth_sendRawTransaction", hexutil.Bytes(content))

	if err != nil {
		return "", err
	}
	return result.String(), nil
}
func (this *Api) CtxGet(hash common.Hash) (*backend.RPCCrossTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	var r *backend.RPCCrossTransaction
	err := this.client.CallContext(ctx, &r, "cross_ctxGet", hash)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (this *Api) GetMakerTx(ctxId common.Hash, contract common.Address, from common.Address, abiData []byte, targetNetworkId *big.Int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	abi, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		return false, err
	}
	input, err := abi.Pack("getMakerTx", ctxId, targetNetworkId)
	msg := simplechain.CallMsg{From: from, To: &contract, Data: input}
	result, err := this.simpleClient.CallContract(ctx, msg, nil)
	if err != nil {
		return false, err
	}
	if new(big.Int).SetBytes(result).Cmp(big.NewInt(0)) > 0 { // error if makerTx is not existed in source-chain
		return true, nil
	} else {
		return false, nil
	}
}
