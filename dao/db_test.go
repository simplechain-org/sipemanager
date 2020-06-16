package dao

import "testing"

var config = &DBConfig{
	Username: "root",
	Password: "root",
	Address:  "localhost",
	Port:     3306,
	DbName:   "sipe",
	Charset:  "utf8mb4",
	MaxIdle:  1000,
	MaxOpen:  2000,
	LogMode:  true,
	Loc:      "Asia/Shanghai",
}

func TestGetDBConnection(t *testing.T) {
	_, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
}
