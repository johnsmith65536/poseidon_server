package redis

import "testing"

func TestGetOnlineUsers(t *testing.T) {
	Init()
	userIds, err := GetUsers()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(userIds)
}

func TestKickAllOnlineUser(t *testing.T) {
	Init()
	err := KickAllUser()
	if err != nil {
		t.Fatal(err)
	}
}