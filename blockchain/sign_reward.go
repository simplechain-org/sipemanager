package blockchain

import (
	"bytes"
	"context"
	"math/big"

	"github.com/simplechain-org/go-simplechain"
	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/common/hexutil"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/simplechain-org/go-simplechain/rlp"
)

type AnchorNodeRewardConfig struct {
	AbiData         []byte
	ContractAddress common.Address
	TargetNetworkId uint64
	AnchorAddress   common.Address
}

//获取奖励池的剩余数量
func (this *Api) GetTotalReward(config *AnchorNodeRewardConfig, callerConfig *CallerConfig) (*big.Int, error) {
	abi, err := abi.JSON(bytes.NewReader(config.AbiData))
	if err != nil {
		return nil, err
	}
	out, err := abi.Pack("getTotalReward", big.NewInt(0).SetUint64(config.TargetNetworkId), config.AnchorAddress)
	if err != nil {
		return nil, err
	}

	msg := simplechain.CallMsg{
		From: callerConfig.From,
		To:   &config.ContractAddress,
		Data: out,
	}
	result, err := this.simpleClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(result), nil
}

//分发交易签名奖励
func (this *Api) AccumulateRewards(config *AnchorNodeRewardConfig, callerConfig *CallerConfig, reward *big.Int) (string, error) {
	abi, err := abi.JSON(bytes.NewReader(config.AbiData))
	if err != nil {
		return "", err
	}
	out, err := abi.Pack("accumulateRewards", big.NewInt(0).SetUint64(config.TargetNetworkId), config.AnchorAddress)
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

//签名数和完成最后一步交易的数量
//获取signCount和finishCount
func (this *Api) GetAnchorWorkCount(config *AnchorNodeRewardConfig, callerConfig *CallerConfig) (*big.Int, *big.Int, error) {
	abi, err := abi.JSON(bytes.NewReader(config.AbiData))
	if err != nil {
		return nil, nil, err
	}
	data, err := abi.Pack("getAnchorWorkCount", big.NewInt(0).SetUint64(config.TargetNetworkId), config.AnchorAddress)
	if err != nil {
		return nil, nil, err
	}

	msg := simplechain.CallMsg{
		From: callerConfig.From,
		To:   &config.ContractAddress,
		Data: data,
	}
	result, err := this.simpleClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, nil, err
	}

	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}

	err = abi.Unpack(out, "getAnchorWorkCount", result)

	if err != nil {
		return nil, nil, err
	}
	return *ret0, *ret1, nil

}

//获取链的单笔签名奖励
func (this *Api) GetChainReward(config *AnchorNodeRewardConfig, callerConfig *CallerConfig) (*big.Int, error) {
	abi, err := abi.JSON(bytes.NewReader(config.AbiData))
	if err != nil {
		return nil, err
	}
	out, err := abi.Pack("getChainReward", big.NewInt(0).SetUint64(config.TargetNetworkId))
	if err != nil {
		return nil, err
	}
	msg := simplechain.CallMsg{
		From: callerConfig.From,
		To:   &config.ContractAddress,
		Data: out,
	}
	result, err := this.simpleClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(result), nil

}

func (this *Api) SetReward(config *AnchorNodeRewardConfig, callerConfig *CallerConfig, reward *big.Int) (string, error) {
	abi, err := abi.JSON(bytes.NewReader(config.AbiData))
	if err != nil {
		return "", err
	}
	out, err := abi.Pack("setReward", big.NewInt(0).SetUint64(config.TargetNetworkId), reward)
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
