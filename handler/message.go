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
	return &thrift.SendMessageResp{ID: msg.Id, CreateTime: msg.CreateTime}, nil
}

func UpdateMessageStatus(ctx context.Context, req *thrift.UpdateMessageStatusReq) (*thrift.UpdateMessageStatusResp, error) {
	err := mysql.UpdateMessageStatus(req.MessageIds,req.UserRelationRequestIds)
	if err != nil {
		return nil, err
	}
	return &thrift.UpdateMessageStatusResp{}, nil
}

func SyncMessage(ctx context.Context, req *thrift.SyncMessageReq) (*thrift.SyncMessageResp, error) {
	messages, err := mysql.SyncMessage(req.UserId, req.MessageId)
	if err != nil {
		return nil, err
	}
	userRelations, err := mysql.SyncUserRelation(req.UserId, req.UserRelationId)
	if err != nil {
		return nil, err
	}
	lastOnlineTime, err := mysql.GetLastOnlineTime(req.UserId)
	if err != nil {
		return nil, err
	}
	return &thrift.SyncMessageResp{Messages: messages, UserRelations: userRelations, LastOnlineTime: lastOnlineTime}, nil
}
