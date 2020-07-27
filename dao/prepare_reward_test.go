package dao

import "testing"

func TestDataBaseAccessObject_GetPrepareRewardCount(t *testing.T) {
	result,err:=obj.GetPrepareRewardCount()
	if err!=nil{
		t.Fatal(err)
		return
	}
	t.Log(result)
}

func TestDataBaseAccessObject_GetPrepareRewardPage(t *testing.T) {
	result,err:=obj.GetPrepareRewardPage(0,10)
	if err!=nil{
		t.Fatal(err)
		return
	}
	t.Log(result)
}
