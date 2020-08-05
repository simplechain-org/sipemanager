package dao

import (
	"fmt"
	"testing"
)

func TestDataBaseAccessObject_WalletExists(t *testing.T) {
	fmt.Println(obj.WalletExists("0xcafc0ec4cb8c123440c3dfbb6cde21240c0c35b8"))
}
