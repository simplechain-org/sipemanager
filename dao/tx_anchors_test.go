package dao

import "testing"

func TestDataBaseAccessObject_QueryAnchors(t *testing.T) {
	config := &DBConfig{
		Username: "root",
		Password: "admin123",
		Address:  "localhost",
		Port:     3306,
		DbName:   "sipe",
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
	result, err := obj.QueryAnchors("202025", "202025", 2, "week")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

}
