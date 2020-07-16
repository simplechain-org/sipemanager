package dao

import "testing"

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
