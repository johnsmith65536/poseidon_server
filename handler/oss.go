package handler

import (
	"context"
	"poseidon/infra/oss"
	"poseidon/thrift"
)

func GetSTSInfo(ctx context.Context, req *thrift.GetSTSInfoReq) (*thrift.GetSTSInfoResp, error) {
	info, err := oss.GetSTSInfo(req.UserId)
	if err != nil {
		return nil, err
	}
	return &thrift.GetSTSInfoResp{SecurityToken: info.SecurityToken, AccessKeyId: info.AccessKeyId, AccessKeySecret: info.AccessKeySecret}, nil
}
