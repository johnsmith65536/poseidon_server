package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
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
		log.Println(http.ListenAndServe(":6060", nil))

	}()
	initHttpServer(ginAddr)
}
