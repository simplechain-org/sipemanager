package dao

import (
	"fmt"
	"testing"
)

func TestDataBaseAccessObject_CreateChain(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)
	chain := &Chain{
		Name:      "主链",
		NetworkId: 1,
		CoinName:  "SIPC",
		Symbol:    "sipc",
	}
	_, err = obj.CreateChain(chain)
	if err != nil {
		t.Fatal(err)
	}
}
func TestDataBaseAccessObject_CreateChain2(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)
	chain := &Chain{
		Name:      "子链",
		NetworkId: 666,
		CoinName:  "GWC",
		Symbol:    "GWC",
	}
	_, err = obj.CreateChain(chain)
	if err != nil {
		t.Fatal(err)
	}
}
func TestDataBaseAccessObject_GetChainIdByContractAddress(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)

	address := "0xBf87aAB36391BE9438819C00A0B6b77Dc665c738"

	chainId, err := obj.GetChainIdByContractAddress(address)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("chainId1=", chainId)

	address = "0xA0392F87E89Fd6959816863c4d0De47BeC38d4C6"

	chainId, err = obj.GetChainIdByContractAddress(address)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("chainId2=", chainId)

}

<<<<<<< HEAD
//func TestDataBaseAccessObject_GetTargetChainId(t *testing.T) {
//	db, err := GetDBConnection(config)
//	if err != nil {
//		t.Fatal(err)
//	}
//	obj := NewDataBaseAccessObject(db)
//
//	var sourceChainId uint = 2
//
//	var targetNetwordId int = 1
//
//	chainId, err := obj.GetTargetChainId(sourceChainId, targetNetwordId)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("chainId=", chainId)
//}
=======
func TestDataBaseAccessObject_GetTargetChainId(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)

	var sourceChainId uint = 2

	var targetNetwordId uint64 = 1

	chainId, err := obj.GetTargetChainId(sourceChainId, targetNetwordId)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("chainId=", chainId)
}
>>>>>>> 手续费报销以及签名奖励
