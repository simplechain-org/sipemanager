package dao

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestDataBaseAccessObject_GetTxTokenList(t *testing.T) {
	config := &DBConfig{
		Username: "root",
		Password: "root",
		Address:  "192.168.3.116",
		Port:     3306,
		DbName:   "sipe_manager",
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
	result, err := obj.GetTxTokenList()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

}

func TestDataBaseAccessObject_CreateChainRegister(t *testing.T) {

	//type ChainRegister struct {
	//	gorm.Model
	//	SourceChainId   uint   `json:"source_chain_id"`
	//	TargetChainId   uint   `json:"target_chain_id"`
	//	Confirm         uint   `json:"confirm"`
	//	AnchorAddresses string `json:"anchor_addresses"`
	//	Status          int    `json:"status"`
	//	StatusText      string `json:"status_text"`
	//	TxHash          string `json:"tx_hash"`
	//	Address         string `json:"address"` // 合约地址
	//}

	chainRegister := &ChainRegister{
		SourceChainId: 1,
		TargetChainId: 2,
		Confirm:       3,
		Status:        1,
		StatusText:    "success",
		TxHash:        "0xijasdiojfaiosjfdsojf",
	}

	id, err := obj.CreateChainRegister(chainRegister)
	if err != nil {
		t.Fatal(err)
		return

	}
	t.Log("id:", id)
}

func TestDataBaseAccessObject_GetChainRegisterPage(t *testing.T) {
	result, err := obj.GetChainRegisterPage(0, 10)
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, o := range result {
		t.Log(o)
	}
}

func TestDataBaseAccessObject_GetChainRegisterCount(t *testing.T) {
	result, err := obj.GetChainRegisterCount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(result)
}

func TestDataBaseAccessObject_GetChainRegister(t *testing.T) {
	chain, err := obj.GetChainRegister(8)
	if err != nil {
		t.Fatal(err)
		return
	}

	type ChainRegisterInfo struct {
		ChainRegisterView
		AnchorNodes []*AnchorNode `json:"anchor_nodes" gorm:"anchor_nodes"`
	}
	chainRegisterInfo := &ChainRegisterInfo{
		ChainRegisterView: *chain,
		AnchorNodes:       make([]*AnchorNode, 0),
	}
	if chain.AnchorAddresses != "" {
		idStrings := strings.Split(chain.AnchorAddresses, ",")
		for _, idStr := range idStrings {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			anchorNode, err := obj.GetAnchorNode(uint(id))
			if err != nil {
				fmt.Println(err)
				continue
			}
			chainRegisterInfo.AnchorNodes = append(chainRegisterInfo.AnchorNodes, anchorNode)
		}

		fmt.Printf("result=%+v\n", chainRegisterInfo)
	}

}
