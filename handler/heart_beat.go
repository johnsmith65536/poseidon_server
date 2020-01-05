package handler

import (
	"context"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"poseidon/thrift"
	"time"
)

func HeartBeat(ctx context.Context, req *thrift.HeartBeatReq) (*thrift.HeartBeatResp, error) {
	err := mysql.UpdateLastOnlineTime(req.UserId, time.Now())
	if err != nil {
		return nil, err
	}
	err = redis.AddUser(req.UserId)
	if err != nil {
		return nil, err
	}
	return &thrift.HeartBeatResp{}, nil
}
