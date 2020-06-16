package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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
	out, err := abi.Pack("chainRegister", big.NewInt(0).SetUint64(config.TargetNetworkId), config.SignConfirmCount, config.AnchorAddresses)

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
	msg := ethereum.CallMsg{
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
