package dao

import "testing"

var config = &DBConfig{
	Username: "root",
	Password: "root",
	Address:  "192.168.3.116",
	Port:     3306,
	DbName:   "sipe_manager_test",
	Charset:  "utf8mb4",
	MaxIdle:  1000,
	MaxOpen:  2000,
	LogMode:  true,
	Loc:      "Asia/Shanghai",
}

var obj *DataBaseAccessObject

func init() {
	db, err := GetDBConnection(config)
	if err != nil {
		panic(err)
	}
	AutoMigrate(db)
	obj = NewDataBaseAccessObject(db)
}

func TestGetDBConnection(t *testing.T) {
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	AutoMigrate(db)
}
