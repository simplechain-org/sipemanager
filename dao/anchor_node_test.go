package dao

import (
	"testing"
)

func TestDataBaseAccessObject_CreateAnchorNode(t *testing.T) {
	//添加锚定节点
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	obj := NewDataBaseAccessObject(db)
	//Name          string `gorm:"name"`    //锚定节点名称
	//Address       string `gorm:"address"` //锚定节点地址
	//SourceChainId uint   `gorm:"source_chain_id"`
	//TargetChainId uint   `gorm:"target_chain_id"`
	//SourceHash    string `gorm:"source_hash"`   //链上的交易哈希
	//TargetHash    string `gorm:"target_hash"`   //链上的交易哈希
	//SourceStatus  uint   `gorm:"source_status"` //链上达成的状态  锚定节点添加成功
	//TargetStatus  uint   `gorm:"target_status"` //链上达成的状态  锚定节点添加成功
	anchorNode := &AnchorNode{
		Name:          "锚定节点1",
		Address:       "0x2d9b3E6b4a446195c912e27c9F3EE592305314ef",
		SourceChainId: 1,
		TargetChainId: 666,
		SourceHash:    "0xhhjjjd",
		TargetHash:    "0xjjkkll",
		SourceStatus:  0,
		TargetStatus:  0,
	}
	id, err := obj.CreateAnchorNode(anchorNode)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("anchor node id:", id)
}

func TestDataBaseAccessObject_CreateAnchorNodeByTx(t *testing.T) {
	//添加锚定节点,开事务
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	obj := NewDataBaseAccessObject(db)
	db2:=obj.BeginTransaction()
	anchorNode := &AnchorNode{
		Name:          "锚定节点2",
		Address:       "0x2d9b3E6b4a446195c912e27c9F3EE592305314eg",
		SourceChainId: 1,
		TargetChainId: 666,
		SourceHash:    "0xhhjjjj",
		TargetHash:    "0xjjkklk",
		SourceStatus:  0,
		TargetStatus:  0,
	}
	id, err := obj.CreateAnchorNodeByTx(db2,anchorNode)
	if err != nil {
		t.Error(err)
		return
	}
	db2.Commit()
	t.Log("CreateAnchorNodeByTx anchor node id:",id)

}

func TestDataBaseAccessObject_UpdateSourceStatus(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	obj := NewDataBaseAccessObject(db)

	err=obj.UpdateSourceStatus(uint(1),uint(1))
	if err!=nil{
		t.Fatal(err)
		return
	}
}

func TestDataBaseAccessObject_UpdateTargetStatus(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	obj := NewDataBaseAccessObject(db)

	err=obj.UpdateTargetStatus(uint(1),uint(1))
	if err!=nil{
		t.Fatal(err)
		return
	}
}

func TestDataBaseAccessObject_GetAnchorNode(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	obj := NewDataBaseAccessObject(db)

	//如果数据不存在，则报错（record not found）
	result,err:=obj.GetAnchorNode(uint(1))
	if err!=nil{
		t.Fatal(err)
		return
	}
	t.Log(result)
}

func TestDataBaseAccessObject_RemoveAnchorNode(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	obj := NewDataBaseAccessObject(db)

	err=obj.RemoveAnchorNode(uint(1))
	if err!=nil{
		t.Fatal(err)
		return
	}
}
