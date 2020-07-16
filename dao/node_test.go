package dao

import (
	"fmt"
	"testing"
)

func TestDataBaseAccessObject_CreateNode(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)
	node := &Node{
		Address: "127.0.0.1",
		Port:    8545,
		ChainId: 1,
		IsHttps: false,
		Name:    "主链节点1",
		UserId:  2,
	}
	id, err := obj.CreateNode(node)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
	node = &Node{
		Address: "127.0.0.1",
		Port:    10546,
		ChainId: 2,
		IsHttps: false,
		Name:    "锚定节点1",
		UserId:  2,
	}
	id, err = obj.CreateNode(node)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

}
func TestDataBaseAccessObject_ListAllNode(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	obj := NewDataBaseAccessObject(db)
	nodes, err := obj.ListAllNode()
	if err == nil {
		for _, node := range nodes {
			fmt.Printf("%+v\n", node)
		}
	}
}
