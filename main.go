package main

import (
	"fmt"
	gitThrift "github.com/apache/thrift/lib/go/thrift"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"poseidon/thrift"
	"time"
)

func initRpcServer(addr string) error {
	handler := &ServiceImpl{}
	processor := thrift.NewServerProcessor(handler)
	serverTransport, err := gitThrift.NewTServerSocket(addr)
	if err != nil {
		return fmt.Errorf("NewTServerSocket failed, err: %+v\n", err)
	}
	transportFactory := gitThrift.NewTBufferedTransportFactory(10000000)
	protocolFactory := gitThrift.NewTBinaryProtocolFactoryDefault()
	server := gitThrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Running at:", addr)
	server.Serve()
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	mysql.Init()
	redis.Init()
	const addr = ":8080"
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

	go initRpcServer(addr)
	initHttpServer(ginAddr)
}
