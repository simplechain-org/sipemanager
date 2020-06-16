package dao

import "testing"

func TestDataBaseAccessObject_CreateUser(t *testing.T) {
	config := &DBConfig{
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
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	//自动迁移仅仅会创建表，缺少列和索引，并且不会改变现有列的类型或删除未使用的列以保护数据。
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Node{})
	db.AutoMigrate(&UserNode{})
	obj := NewDataBaseAccessObject(db)
	user := &User{
		Username: "yangdamin",
		Password: "123456",
	}
	id, err := obj.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

}

func TestDataBaseAccessObject_UserIsValid(t *testing.T) {
	config := &DBConfig{
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
	db, err := GetDBConnection(config)
	if err != nil {
		t.Fatal(err)
	}
	//自动迁移仅仅会创建表，缺少列和索引，并且不会改变现有列的类型或删除未使用的列以保护数据。
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Node{})
	db.AutoMigrate(&UserNode{})
	obj := NewDataBaseAccessObject(db)
	t.Log(obj.UserIsValid("yangdamin", "123456"))
	t.Log(obj.UserIsValid("yangdamin", "1234567"))
}
