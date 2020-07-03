package blockchain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type Order struct {
	Value            *big.Int
	TxId             common.Hash
	TxHash           common.Hash
	From             common.Address
	To               common.Address
	BlockHash        common.Hash
	DestinationValue *big.Int
	Data             []byte
	V                []*big.Int
	R                [][32]byte
	S                [][32]byte
}

type HandleOrder struct {
	TxHash      string
	SourceValue *big.Int
	TargetValue *big.Int
	NetworkId   uint64
	CtxId       string
}

//测试通过
//要注意合约的地址是否设置正确，数据库表contracts的address字段
func (this *Api) Consume(foundCtxId string, contractConfig *ContractConfig, callerConfig *CallerConfig) (*HandleOrder, error) {
	var signatures map[string]map[uint64][]*RPCCrossTransaction
	err := this.client.CallContext(context.Background(), &signatures, "eth_ctxContent")
	if err != nil {
		return nil, err
	}
	abi, err := abi.JSON(bytes.NewReader(contractConfig.AbiData))
	if err != nil {
		return nil, err
	}
	for remoteId, value := range signatures["remote"] {
		for _, v := range value {
			fmt.Println("v.CTxId.String() =", v.CTxId.String())
			if foundCtxId == v.CTxId.String() {
				r := make([][32]byte, 0, len(v.R))
				s := make([][32]byte, 0, len(v.S))
				vv := make([]*big.Int, 0, len(v.V))
				for i := 0; i < len(v.R); i++ {
					rone := common.LeftPadBytes(v.R[i].ToInt().Bytes(), 32)
					var a [32]byte
					copy(a[:], rone)
					r = append(r, a)
					sone := common.LeftPadBytes(v.S[i].ToInt().Bytes(), 32)
					var b [32]byte
					copy(b[:], sone)
					s = append(s, b)
					vv = append(vv, v.V[i].ToInt())
				}
				chainId := big.NewInt(int64(remoteId))
				var ord Order
				ord.Value = v.Value.ToInt()
				ord.TxId = v.CTxId
				ord.TxHash = v.TxHash
				ord.From = v.From
				ord.To = v.To
				ord.BlockHash = v.BlockHash
				ord.DestinationValue = v.DestinationValue.ToInt()
				ord.Data = v.Input
				ord.V = vv
				ord.R = r
				ord.S = s
				out, err := abi.Pack("taker", &ord, chainId)
				if err != nil {
					return nil, err
				}
				nonce, err := this.simpleClient.PendingNonceAt(context.Background(), callerConfig.From)
				if err != nil {
					return nil, err
				}
				gasPrice, err := this.simpleClient.SuggestGasPrice(context.Background())

				if err != nil {
					return nil, err
				}
				//必须给足了值
				sourceValue := v.DestinationValue.ToInt()
				msg := ethereum.CallMsg{
					From:     callerConfig.From,
					To:       &contractConfig.ContractAddress,
					Data:     out,
					Value:    sourceValue,
					GasPrice: gasPrice,
				}
				gasLimit, err := this.simpleClient.EstimateGas(context.Background(), msg)
				if err != nil {
					return nil, err
				}

				transaction := types.NewTransaction(nonce, contractConfig.ContractAddress, sourceValue, gasLimit, gasPrice, out)

				transaction, err = types.SignTx(transaction, types.NewEIP155Signer(big.NewInt(0).SetInt64(int64(callerConfig.NetworkId))), callerConfig.PrivateKey)
				if err != nil {
					return nil, err
				}

				content, err := rlp.EncodeToBytes(transaction)
				if err != nil {
					return nil, err
				}
				var result common.Hash

				err = this.client.CallContext(context.Background(), &result, "eth_sendRawTransaction", hexutil.Bytes(content))

				if err != nil {
					return nil, err
				}
				//接单的时候需要把字段反过来
				handleOrder := &HandleOrder{
					TxHash:      result.Hex(),
					TargetValue: v.DestinationValue.ToInt(),
					SourceValue: v.Value.ToInt(),
					NetworkId:   remoteId,
					CtxId:       foundCtxId,
				}

				return handleOrder, nil
			}
		}

	}
	return nil, errors.New(fmt.Sprintf("no transaction match for %s", foundCtxId))
}
