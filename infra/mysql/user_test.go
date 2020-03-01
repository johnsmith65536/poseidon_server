package mysql

import "testing"

func  TestSearchUser(t *testing.T) {
	Init()
	users,err := SearchUser("te")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(users)
}
