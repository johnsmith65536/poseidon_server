package main

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

func initRpcServer(addr string) error {
	handler := &ServiceImpl{}
	processor := NewServerProcessor(handler)
	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return fmt.Errorf("NewTServerSocket failed, err: %+v\n", err)
	}
	transportFactory := thrift.NewTBufferedTransportFactory(10000000)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Running at:", addr)
	server.Serve()
	return nil
}

func main() {
	const addr = ":8080"
	initRpcServer(addr)
}
