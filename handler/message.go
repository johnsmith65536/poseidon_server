package handler

import (
	"context"
	"poseidon/entity"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"poseidon/thrift"
	"strconv"
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
	msg, err := mysql.WriteMessage(req.UserIdSend, req.IdRecv, 0, req.Content, time.Now(), req.ContentType, req.MessageType, false)
	if err != nil {
		return nil, err
	}
	broadcastMsg := map[string]interface{}{
		"Id":          msg.Id,
		"UserIdSend":  msg.UserIdSend,
		"UserIdRecv":  msg.UserIdRecv,
		"GroupId":     msg.GroupId,
		"Content":     msg.Content,
		"CreateTime":  msg.CreateTime,
		"ContentType": msg.ContentType,
		"MsgType":     msg.MsgType,
	}
	switch req.ContentType {
	case int32(entity.Text):
		;
	case int32(entity.ObjectData):
		objId, err := strconv.ParseInt(req.Content, 10, 64)
		if err != nil {
			return nil, err
		}
		object, err := mysql.GetObject(objId)
		if err != nil {
			return nil, err
		}
		broadcastMsg["ObjectETag"] = object.ETag
		broadcastMsg["ObjectName"] = object.Name
	}
	go redis.BroadcastMessage(req.IdRecv, broadcastMsg, redis.Chat)
	return &thrift.SendMessageResp{ID: msg.Id, CreateTime: msg.CreateTime}, nil
}

func UpdateMessageStatus(ctx context.Context, req *thrift.UpdateMessageStatusReq) (*thrift.UpdateMessageStatusResp, error) {
	err := mysql.UpdateMessageStatus(req.MessageIds, req.UserRelationRequestIds)
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
	var objIds []int64
	for _, message := range messages {
		if message.ContentType == int32(entity.ObjectData) {
			objId, err := strconv.ParseInt(message.Content, 10, 64)
			if err != nil {
				return nil, err
			}
			objIds = append(objIds, objId)
		}
	}
	objects, err := mysql.SyncObject(objIds)
	if err != nil {
		return nil, err
	}
	lastOnlineTime, err := mysql.GetLastOnlineTime(req.UserId)
	if err != nil {
		return nil, err
	}
	return &thrift.SyncMessageResp{Messages: messages, UserRelations: userRelations, Objects: objects, LastOnlineTime: lastOnlineTime}, nil
}

func FetchMessageStatus(ctx context.Context, req *thrift.FetchMessageStatusReq) (*thrift.FetchMessageStatusResp, error) {
	messageStatus, err := mysql.GetMessageStatus(req.MessageIds)
	if err != nil {
		return nil, err
	}
	relationStatus, err := mysql.GetRelationStatus(req.UserRelationRequestIds)
	if err != nil {
		return nil, err
	}
	return &thrift.FetchMessageStatusResp{MessageIds: messageStatus, UserRelationRequestIds: relationStatus}, nil
}
