package mysql

import (
	_ "github.com/go-sql-driver/mysql"
)

/*func TestFetchOfflineUserRelation(t *testing.T) {
	Init()
	db.LogMode(true)
	resp,err := FetchOfflineUserRelation(1709,-1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
*/