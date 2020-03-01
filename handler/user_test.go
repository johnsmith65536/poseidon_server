package handler

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"poseidon/infra/mysql"
	"poseidon/thrift"
	"testing"
)

func TestSearchUser(t *testing.T) {
	mysql.Init()
	ctx := context.Background()
	req := thrift.SearchUserReq{UserId: 1566, Data: ""}
	resp, err := SearchUser(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*resp)
}
