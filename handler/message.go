package handler

import (
	"context"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"poseidon/thrift"
	"time"
)

//Id         int64
//UserIdSend int64
//UserIdRecv int64
//GroupId    int64
//Content    string
//CreateTime int64
//MsgType    int32
//Delivered  bool
//Read       bool

func SendMessage(ctx context.Context, req *thrift.SendMessageReq) (*thrift.SendMessageResp, error) {
	msg, err := mysql.WriteMessage(req.UserIdSend, req.IdRecv, 0, req.Content, time.Now(), req.MessageType, false)
	if err != nil {
		return nil, err
	}
	go redis.BroadcastMessage(req.IdRecv, map[string]interface{}{
		"Id":         msg.Id,
		"UserIdSend": msg.UserIdSend,
		"UserIdRecv": msg.UserIdRecv,
		"GroupId":    msg.GroupId,
		"Content":    msg.Content,
		"CreateTime": msg.CreateTime,
		"MsgType":    msg.MsgType,
	}, redis.Chat)
	return &thrift.SendMessageResp{ID: msg.Id}, nil
}

func MessageDelivered(ctx context.Context, req *thrift.MessageDeliveredReq) (*thrift.MessageDeliveredResp, error) {
	err := mysql.MessageDelivered(req.MsgId)
	if err != nil {
		return nil, err
	}
	return &thrift.MessageDeliveredResp{Status: 0}, nil
}

func FetchOfflineMessage(ctx context.Context, req *thrift.FetchOfflineMessageReq) (*thrift.FetchOfflineMessageResp, error) {
	messages, err := mysql.FetchOfflineMessage(req.UserId, req.MessageId)
	if err != nil {
		return nil, err
	}
	userRelations, err := mysql.FetchOfflineUserRelation(req.UserId, req.UserRelationId)
	if err != nil {
		return nil, err
	}
	return &thrift.FetchOfflineMessageResp{Messages: messages, UserRelations: userRelations}, nil
}
