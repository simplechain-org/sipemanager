package dao

import (
	"context"
	_ "errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	_ "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"math/big"
	_ "net/url"
	_ "regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetTxByHash(t *testing.T) {
	config := &DBConfig{
		Username: "root",
		Password: "admin123",
		Address:  "localhost",
		Port:     3306,
		DbName:   "sipe",
		Charset:  "utf8mb4",
		MaxIdle:  1000,
		MaxOpen:  2000,
		LogMode:  true,
		Loc:      "Asia/Shanghai",
	}
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)
	id, err := obj.GetTxByHash("0xa8c74f5e603fa7b541bd032dd392af4a671bdab39152f68ebbf69aabdd1785d9")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

}

func TestGetTxReplace(t *testing.T) {
	config := &DBConfig{
		Username: "root",
		Password: "admin123",
		Address:  "localhost",
		Port:     3306,
		DbName:   "sipe",
		Charset:  "utf8mb4",
		MaxIdle:  1000,
		MaxOpen:  2000,
		LogMode:  true,
		Loc:      "Asia/Shanghai",
	}
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)
	urlStr := "http://101.68.74.172:8556"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	client, err := rpc.DialContext(ctx, urlStr)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	ethClient := ethclient.NewClient(client)

	block, err := ethClient.BlockByNumber(ctx, big.NewInt(0).SetInt64(114279))
	if err != nil {
		t.Fatal(err)
	}
	if len(block.Transactions()) > 0 {
		for index, transaction := range block.Transactions() {
			netId, err := ethClient.NetworkID(ctx)
			// 获取from地址
			msg, err := transaction.AsMessage(types.NewEIP155Signer(netId))
			if err != nil {
				logrus.Error(err.Error())
			}
			from := msg.From().Hex()
			receipt, err := ethClient.TransactionReceipt(ctx, transaction.Hash())
			if err != nil {
				logrus.Error("TransactionReceipt:", err)
			}
			// hexutil.Encode(transaction.Data()),
			txRecord := Transaction{
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
				ChainId:          2,
			}
			//var encodeInput string
			if transaction.To() != nil {
				txRecord.To = strings.ToLower(transaction.To().Hex())
			}
			fmt.Printf("%+v\n", txRecord)

			//encodeInput = hexutil.Encode(transaction.Data())
			contract, err := obj.GetContractById(1)
			abiParsed, err := abi.JSON(strings.NewReader(contract.Abi))

			fmt.Printf("---- %+v\n", abiParsed)
			sigdata := transaction.Data()[:4]

			argdata := transaction.Data()[4:]
			fmt.Println(argdata)
			method, err := abiParsed.MethodById(sigdata)
			var args MakerFinish
			if err := method.Inputs.Unpack(&args, argdata); err != nil {
				fmt.Println("UnpackValues err=", err)
			}
			if err != nil {
				fmt.Println("abi.MethodById err=", err)
			}
			t.Log(args.Rtx.TxId.Hex())
		}
	}

}

//func TestDataBaseAccessObject_QueryTxByHours(t *testing.T) {
//	config := &DBConfig{
//		Username: "root",
//		Password: "admin123",
//		Address:  "localhost",
//		Port:     3306,
//		DbName:   "sipe",
//		Charset:  "utf8mb4",
//		MaxIdle:  1000,
//		MaxOpen:  2000,
//		LogMode:  true,
//		Loc:      "Asia/Shanghai",
//	}
//	db, err := GetDBConnection(config)
//	if err != nil {
//		t.Fatal(err)
//	}
//	obj := NewDataBaseAccessObject(db)
//	qerr := obj.QueryTxByHours("0x17529b05513e5595ceff7f4fb1e06512c271a540", "makerFinish", 2, 1)
//	if qerr != nil {
//		t.Fatal(qerr)
//	}
//	//t.Log(id)
//
//}
