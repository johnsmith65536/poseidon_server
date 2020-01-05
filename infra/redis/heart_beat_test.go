package redis

import (
	"fmt"
	"testing"
	"time"
)

func TestHeartBeat(t *testing.T) {
	Init()
	go func() {
		t.Log("1 1")
		for {
			err := HeartBeat()
			t.Log("HeartBeat 1")
			if err != nil {
				fmt.Println(err)
				return
			}
			t.Log("HeartBeat success")
			time.Sleep(time.Second * 2)
		}
	}()
	time.Sleep(time.Hour)
}
