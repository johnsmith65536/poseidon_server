package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
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

func TestGetRelationStatus(t *testing.T) {
	Init()
	res, err := GetRelationStatus([]int64{1,62,63})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
