package redis

import (
	"github.com/go-redis/redis"
)

const Addr = "localhost:6379"

var redisCli *redis.Client

func Init() {
	redisCli = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: "",
		DB:       0,
	})
}

