package dao

import "testing"


func TestDataBaseAccessObject_GetContractPage(t *testing.T) {
	result, err := obj.GetContractPage(0, 10,"deployed")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("deployed=",len(result))
	result, err = obj.GetContractPage(0, 10,"not_deployed")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("not_deployed=",len(result))
	result, err = obj.GetContractPage(10, 10,"all")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("all=",len(result))
}

func TestDataBaseAccessObject_GetContractCount(t *testing.T) {
	result, err := obj.GetContractCount("deployed")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("deployed=", result)
	result, err = obj.GetContractCount("not_deployed")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("not_deployed=", result)
	result, err = obj.GetContractCount("all")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("all=", result)
}

func TestDataBaseAccessObject_GetContractCount1(t *testing.T) {
	result, err := obj.GetContractPage(10, 10,"all")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("all=",len(result))
	count, err := obj.GetContractCount("all")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("count=", count)
}
