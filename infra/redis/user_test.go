package redis

import "testing"

func TestGetUsers(t *testing.T) {
	Init()
	userIds, err := GetUsers()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(userIds)
}

//func TestKickAllUser(t *testing.T) {
//	Init()
//	err := KickAllUser()
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func TestAddUser(t *testing.T) {
	Init()
	err := AddUser(1)
	if err != nil {
		t.Fatal(err)
	}
}
