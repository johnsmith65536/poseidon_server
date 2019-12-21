package handler

import (
	"context"
	"poseidon/infra/mysql"
	"poseidon/thrift"
	"poseidon/utils"
	"time"
)

func CreateUser(ctx context.Context, req *thrift.CreateUserReq) (*thrift.CreateUserResp, error) {
	userId := utils.GenerateId(4)
	err := mysql.CreateUser(userId, req.Password, req.NickName, time.Now())
	if err != nil {
		return nil, err
	}
	return &thrift.CreateUserResp{UserId: userId}, nil
}
