package dao

import (
	"encoding/json"
	"testing"
)

func TestDataBaseAccessObject_CreateAnchorNode(t *testing.T) {
	//添加锚定节点
	anchorNode := &AnchorNode{
		Name:          "锚定节点4",
		Address:       "0x2d9b3E6b4a446195c912e27c9F3EE592305314ef",
		SourceChainId: 1,
		TargetChainId: 666,
		SourceHash:    "0xcc8871f8f6b536e06b4465ab1635a28a307301d251e30475e68608150bb38987",
		TargetHash:    "0xac2871f8f6b536e06b4465ab1635a28a307301d251e30475e68608150bb38985",
		SourceStatus:  1,
		TargetStatus:  1,
	}
	id, err := obj.CreateAnchorNode(anchorNode)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("anchor node id:", id)


	//添加锚定节点
	anchorNode = &AnchorNode{
		Name:          "锚定节点5",
		Address:       "0x3d9b3E6b4a446195c912e27c9F3EE592305314ef",
		SourceChainId: 1,
		TargetChainId: 666,
		SourceHash:    "0xcc8871f8f6b536e06b4465ab1635a28a307301d251e30475e68608150bb38987",
		TargetHash:    "0xac2871f8f6b536e06b4465ab1635a28a307301d251e30475e68608150bb38985",
		SourceStatus:  1,
		TargetStatus:  1,
	}
	id, err = obj.CreateAnchorNode(anchorNode)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("anchor node id:", id)


	//添加锚定节点
	anchorNode = &AnchorNode{
		Name:          "锚定节点6",
		Address:       "0x4d9b3E6b4a446195c912e27c9F3EE592305314ef",
		SourceChainId: 1,
		TargetChainId: 666,
		SourceHash:    "0xcc8871f8f6b536e06b4465ab1635a28a307301d251e30475e68608150bb38987",
		TargetHash:    "0xac2871f8f6b536e06b4465ab1635a28a307301d251e30475e68608150bb38985",
		SourceStatus:  1,
		TargetStatus:  1,
	}
	id, err = obj.CreateAnchorNode(anchorNode)
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
func TestDataBaseAccessObject_GetAnchorNodePage(t *testing.T) {
	result,err:=obj.GetAnchorNodePage(0,10,0)
	if err!=nil{
		t.Error(err)
		return
	}
	data,err:=json.Marshal(result)
	if err!=nil{
		t.Error(err)
		return
	}
	t.Log(string(data))
}
