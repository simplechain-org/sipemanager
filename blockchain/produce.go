package blockchain

import (
	"bytes"
	"context"
	"github.com/simplechain-org/go-simplechain/common/hexutil"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/simplechain-org/go-simplechain/rlp"
	"math/big"

	"github.com/simplechain-org/go-simplechain"
	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"github.com/simplechain-org/go-simplechain/common"
)

type ContractConfig struct {
	AbiData         []byte
	ContractAddress common.Address
	TargetChainId   uint64
}

type ChangeParam struct {
	SourceValue *big.Int //付出
	TargetValue *big.Int //得到
	Input       []byte   //附带数据
}

func (this *Api) Produce(contractConfig *ContractConfig, changeParam *ChangeParam, callerConfig *CallerConfig) (string, error) {
	abi, err := abi.JSON(bytes.NewReader(contractConfig.AbiData))
	if err != nil {
		return "", err
	}
	remoteChainId := big.NewInt(0).SetUint64(contractConfig.TargetChainId)
	out, err := abi.Pack("makerStart", remoteChainId, changeParam.TargetValue,common.HexToAddress("0x0"), changeParam.Input)
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
		To:       &contractConfig.ContractAddress,
		Data:     out,
		Value:    changeParam.SourceValue,
		GasPrice: gasPrice,
	}
	gasLimit, err := this.simpleClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", err
	}
	transaction := types.NewTransaction(nonce, contractConfig.ContractAddress, changeParam.SourceValue, gasLimit, gasPrice, out)

	transaction, err = types.SignTx(transaction,
		types.NewEIP155Signer(big.NewInt(0).SetInt64(int64(callerConfig.NetworkId))), callerConfig.PrivateKey)

	content, err := rlp.EncodeToBytes(transaction)

	if err != nil {
		return "", err
	}

	var result common.Hash

	err = this.client.CallContext(context.Background(), &result, "eth_sendRawTransaction", hexutil.Bytes(content))

	if err != nil {
		return "", err
	}
	return result.Hex(), nil
}
