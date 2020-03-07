package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	mysql.Init()
	redis.Init()
	const ginAddr = ":8081"

	go func() {
		for {
			err := redis.HeartBeat()
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Second * 5)
		}
	}()

	initHttpServer(ginAddr)
}
