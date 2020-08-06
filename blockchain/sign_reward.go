package blockchain

import (
	"bytes"
	"context"
	"fmt"
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
//ok
func (this *Api) GetTotalReward(config *AnchorNodeRewardConfig, callerConfig *CallerConfig) (*big.Int, error) {
	abi, err := abi.JSON(bytes.NewReader(config.AbiData))
	if err != nil {
		return nil, fmt.Errorf("读取合约的abi文件json内容出错:%s", err.Error())
	}
	out, err := abi.Pack("getTotalReward", big.NewInt(0).SetUint64(config.TargetNetworkId))
	if err != nil {
		return nil, fmt.Errorf("getTotalReward合约方法打包错误：%s", err.Error())
	}
	msg := simplechain.CallMsg{
		From: callerConfig.From,
		To:   &config.ContractAddress,
		Data: out,
	}
	result, err := this.simpleClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("调用合约方法出错:%s", err.Error())
	}
	return big.NewInt(0).SetBytes(result), nil
}

//发放交易签名奖励
func (this *Api) AccumulateRewards(config *AnchorNodeRewardConfig, callerConfig *CallerConfig, reward *big.Int) (string, error) {

	abi, err := abi.JSON(bytes.NewReader(config.AbiData))
	if err != nil {
		return "", err
	}
	out, err := abi.Pack("accumulateRewards", big.NewInt(0).SetUint64(config.TargetNetworkId), config.AnchorAddress, reward)
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
		fmt.Println("GetAnchorWorkCount abi.JSON", err.Error())
		return nil, nil, err
	}
	data, err := abi.Pack("getAnchorWorkCount", big.NewInt(0).SetUint64(config.TargetNetworkId), config.AnchorAddress)
	if err != nil {
		fmt.Println("GetAnchorWorkCount abi.Pack", err.Error())
		return nil, nil, err
	}

	msg := simplechain.CallMsg{
		From: callerConfig.From,
		To:   &config.ContractAddress,
		Data: data,
	}
	result, err := this.simpleClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		fmt.Println("GetAnchorWorkCount CallContract", err.Error())
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
	if len(result) == 0 {
		return big.NewInt(0), big.NewInt(0), nil
	}
	err = abi.Unpack(out, "getAnchorWorkCount", result)
	if err != nil {
		fmt.Println("GetAnchorWorkCount Unpack", err.Error())
		return nil, nil, err
	}
	return *ret0, *ret1, nil

}

//获取链的单笔预收手续费
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

type SetRewardConfig struct {
	AbiData         []byte
	ContractAddress common.Address
	TargetNetworkId uint64
	AnchorAddress   common.Address
}

//设置链的单笔预收手续费
func (this *Api) SetReward(config *SetRewardConfig, callerConfig *CallerConfig, reward *big.Int) (string, error) {
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
