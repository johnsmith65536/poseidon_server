package handler

import (
	"context"
	"poseidon/infra/mysql"
	"poseidon/thrift"
)

func CreateObject(ctx context.Context, req *thrift.CreateObjectReq) (*thrift.CreateObjectResp, error) {
	id, err := mysql.CreateObject(req.ETag, req.Name)
	if err != nil {
		return nil, err
	}
	return &thrift.CreateObjectResp{ID: id}, nil
}
