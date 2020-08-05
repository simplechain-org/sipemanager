package dao

import (
	"fmt"
	"testing"
)

func TestDataBaseAccessObject_CreateNode(t *testing.T) {
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
	nodes, err := obj.ListAllNode()
	if err == nil {
		for _, node := range nodes {
			fmt.Printf("%+v\n", node)
		}
	}
}

func TestDataBaseAccessObject_ListNodeByUserId(t *testing.T) {
	result, err := obj.ListNodeByUserId(1)
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, o := range result {
		t.Log(o)
	}

}
func TestDataBaseAccessObject_ChainIdExists(t *testing.T) {
	fmt.Println("chain", obj.ChainIdExists(5))
}

func TestDataBaseAccessObject_ListNodeByUserIdPage(t *testing.T) {
	result, err := obj.ListNodeByUserIdPage(0, 10, 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(result))
}
func TestDataBaseAccessObject_GetNodeByUserIdCount(t *testing.T) {
	result, err := obj.GetNodeByUserIdCount(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
