package mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSyncMessage(t *testing.T) {
	Init()
	msgs, err := SyncMessage(1709, -1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msgs)
}

func TestUpdateMessageStatus(t *testing.T) {
	Init()
	err := UpdateMessageStatus(map[int64]int32{1:1},map[int64]int32{})
	if err != nil {
		t.Fatal(err)
	}
}