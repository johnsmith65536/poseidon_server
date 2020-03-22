package mysql

import "testing"

func TestGetRelation(t *testing.T) {
	Init()
	friendUserIds, err := GetRelation(1566, []int64{1566, 7803, 3123, 534, 3123, 5433, 7803})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(friendUserIds)
}

func TestDeleteFriend(t *testing.T) {
	Init()
	err := DeleteFriend(1, 2)
	if err != nil {
		t.Fatal(err)
	}
}
