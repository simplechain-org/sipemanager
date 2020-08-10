package dao

import (
	"testing"
)

func TestDataBaseAccessObject_CreateSignRewardLog(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)
	signRewardLog := &SignRewardLog{
		AnchorNodeId:    1,
		AnchorNodeName:  "锚定节点1",
		TotalReward:     "100",
		Rate:            "20%",
		Reward:          "10",
		TransactionHash: "交易哈希值1",
		Coin:            "SIPC",
		Sender:          "0xddddfffffsafdsfadsf",
		Status:          1,
		FinishCount:     100,
		SignCount:       20,
	}
	id, err := obj.CreateSignRewardLog(signRewardLog)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}

func TestDataBaseAccessObject_UpdateSignRewardLogStatus(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)
	err = obj.UpdateSignRewardLogStatus(uint(1), 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDataBaseAccessObject_GetSignRewardLogPage(t *testing.T) {
	result, err := obj.GetSignRewardLogPage(10, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(result))
}

func TestDataBaseAccessObject_GetSignRewardLogCount(t *testing.T) {
	result, err := obj.GetSignRewardLogCount(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetSignRewardLogCount:", result)
}
