package handler

import (
	"context"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"poseidon/thrift"
)

func Login(ctx context.Context, req *thrift.LoginReq) (*thrift.LoginResp, error) {
	ok, err := mysql.Login(req.UserId, req.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &thrift.LoginResp{Success: ok}, nil
	}
	err = redis.AddUser(req.UserId)
	if err != nil {
		return nil, err
	}
	return &thrift.LoginResp{Success: ok}, err
}

func Logout(ctx context.Context, req *thrift.LogoutReq) (*thrift.LogoutResp, error) {
	err := redis.KickUser(req.UserId)
	if err != nil {
		return nil, err
	}
	return &thrift.LogoutResp{}, err
}
