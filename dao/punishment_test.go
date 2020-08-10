package dao

import "testing"

func TestDataBaseAccessObject_GetPunishmentPage(t *testing.T) {
	result, err := obj.GetPunishmentPage(10, 10, 0)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(len(result))
}

func TestDataBaseAccessObject_GetPunishmentCount(t *testing.T) {
	result, err := obj.GetPunishmentCount(0)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(result)
}
