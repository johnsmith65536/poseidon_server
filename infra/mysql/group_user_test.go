package mysql

import "testing"

func TestGetLastReadMsgId(t *testing.T) {
	Init()
	ret, err := GetLastReadMsgId(5012)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}

func TestUpdateLastReadMsgId(t *testing.T) {
	Init()
	data := map[int64]int64{4250: -1}
	err := UpdateLastReadMsgId(7803, data)
	if err != nil {
		t.Fatal(err)
	}
}
