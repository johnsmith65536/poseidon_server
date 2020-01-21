package mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestFetchOfflineMessage(t *testing.T) {
	Init()
	msgs, err := FetchOfflineMessage(1709, -1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msgs)
}
