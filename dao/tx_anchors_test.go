package dao

import "testing"

func TestDataBaseAccessObject_QueryAnchors(t *testing.T) {
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
	tokenList, err := obj.GetTxTokenList()
	token := tokenList["1,2"]
	result, err := obj.TokenListAnchorCount(token, "202028", "202028", "week", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

}
