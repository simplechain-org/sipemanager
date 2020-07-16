package dao

import (
	"testing"
)

func TestDataBaseAccessObject_CreateServiceChargeLog(t *testing.T) {
	//CreateServiceChargeLog
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	obj := NewDataBaseAccessObject(db)

	//type ServiceChargeLog struct {
	//	gorm.Model
	//	AnchorNodeId    uint     `gorm:"anchor_node_id"`   //锚定节点编号
	//	AnchorNodeName  string   `gorm:"anchor_node_name"` //锚定节点名称，冗余方便查询
	//	TransactionHash string   `gorm:"transaction_hash"` //交易哈希
	//	Fee             *big.Int `gorm:"fee"`              //报销手续费
	//	Coin            string   `gorm:"coin"`             //报销的币种
	//	Sender          string   `gorm:"sender"`           //出账账户地址
	//	Status          uint     `gorm:"status"`           //状态
	//}
	serviceChargeLog := &ServiceChargeLog{
		AnchorNodeId:    1,
		TransactionHash: "交易哈希",
		Fee:             "1000000000000000000000000",
		Coin:            "SIPC",
		Sender:          "0xisdjfiasdjfidosjds",
		Status:          1,
	}
	id, err := obj.CreateServiceChargeLog(serviceChargeLog)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("CreateServiceChargeLog:", id)

}

func TestDataBaseAccessObject_UpdateServiceChargeLogSourceStatus(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	obj := NewDataBaseAccessObject(db)
	err = obj.UpdateServiceChargeLogSourceStatus(uint(1), 1)
	if err != nil {
		t.Fatal(err)
		return
	}
}

//
func TestDataBaseAccessObject_GetServiceChargeLogCount(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	obj := NewDataBaseAccessObject(db)
	count, err := obj.GetServiceChargeLogCount(uint(1))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetServiceChargeLogCount:", count)
}

func TestDataBaseAccessObject_GetServiceChargeLogPage(t *testing.T) {
	result, err := obj.GetServiceChargeLogPage(0, 10, uint(1))
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, o := range result {
		t.Log("GetServiceChargeLogPage:", o)
	}

}

func TestDataBaseAccessObject_GetServiceChargeLog(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	obj := NewDataBaseAccessObject(db)
	serviceChargeLog, err := obj.GetServiceChargeLog(uint(1))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("serviceChargeLog", serviceChargeLog)
}

func TestDataBaseAccessObject_GetServiceChargeSumFee(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	obj := NewDataBaseAccessObject(db)
	sum, err := obj.GetServiceChargeSumFee(uint(1), "SIPC")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetServiceChargeSumFee:", sum)
}
