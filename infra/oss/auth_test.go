package oss

import (
	"math/rand"
	"poseidon/infra/mysql"
	"testing"
	"time"
)

func TestGetSTSInfo(t *testing.T) {
	mysql.Init()
	rand.Seed(time.Now().UnixNano())
	res,err := GetSTSInfo(100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*res)
}