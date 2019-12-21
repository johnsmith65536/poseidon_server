package handler

import (
	"context"
	"fmt"
	"poseidon/thrift"
)

func HeartBeat(ctx context.Context, req *thrift.HeartBeatReq) (resp *thrift.HeartBeatResp, err error) {
	fmt.Printf("%+v\n", req)
	return &thrift.HeartBeatResp{}, nil
}