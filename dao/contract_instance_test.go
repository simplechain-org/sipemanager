package dao

import "testing"

func TestDataBaseAccessObject_GetContractInstanceCount(t *testing.T) {
	result, err := obj.GetContractInstancePage(20, 10)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("result", result)
	count, err := obj.GetContractInstanceCount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("count", count)
}
