package dao

import "testing"

func TestDataBaseAccessObject_GetPunishmentPage(t *testing.T) {
	result, err := obj.GetPunishmentPage(0, 10, 1)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(result)
}
