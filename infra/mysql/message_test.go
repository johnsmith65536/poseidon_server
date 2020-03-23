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

func TestGetGroupLastMsgId(t *testing.T) {
	Init()
	id, err := GetGroupLastMsgId(9208)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}
