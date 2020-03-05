package mysql

import "testing"

func TestCreateObject(t *testing.T) {
	Init()
	id, err := CreateObject("e_tag", "nameaa")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}
