package dao

import (
	"fmt"
	"testing"
)

func TestDataBaseAccessObject_WalletExists(t *testing.T) {
	fmt.Println(obj.WalletExists("0xcafc0ec4cb8c123440c3dfbb6cde21240c0c35b8"))
}

func TestDataBaseAccessObject_GetWalletViewCount(t *testing.T) {
	result,err:=obj.GetWalletViewCount(1)
	if err!=nil{
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestDataBaseAccessObject_GetWalletViewPage(t *testing.T) {
	result,err:=obj.GetWalletViewPage(1,0,10)
	if err!=nil{
		t.Error(err)
		return
	}
	t.Log(len(result))
}
