package dao

import "testing"

func TestDataBaseAccessObject_CreateUser(t *testing.T) {
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
	t.Log(obj.UserIsValid("yangdamin", "123456"))
	t.Log(obj.UserIsValid("yangdamin", "1234567"))
}
