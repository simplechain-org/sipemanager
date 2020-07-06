package blockchain

import (
	"context"
	"crypto/ecdsa"
	"github.com/simplechain-org/go-simplechain"
	"github.com/simplechain-org/go-simplechain/accounts/keystore"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/common/hexutil"
	"github.com/simplechain-org/go-simplechain/core/types"
	"github.com/simplechain-org/go-simplechain/rlp"
	"math/big"
)

func GetPrivateKey(keyjson []byte, passphrase string) (*ecdsa.PrivateKey, error) {
	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return nil, err
	}
	return key.PrivateKey, nil
}

func (this *Api) DeployContract(from common.Address, amount *big.Int, data []byte, chainId uint64, privateKey *ecdsa.PrivateKey) (string, error) {
	nonce, err := this.simpleClient.PendingNonceAt(context.Background(), from)
	if err != nil {
		return "", err
	}
	gasPrice, err := this.simpleClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	msg := simplechain.CallMsg{
		From:     from,
		GasPrice: gasPrice,
		Data:     data,
		Value:    amount,
	}
	gasLimit, err := this.simpleClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	transaction := types.NewContractCreation(nonce, amount, gasLimit+90000, gasPrice, data)

	transaction, err = types.SignTx(transaction, types.NewEIP155Signer(big.NewInt(0).SetInt64(int64(chainId))), privateKey)

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
