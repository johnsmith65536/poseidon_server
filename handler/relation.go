package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"poseidon/entity"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"poseidon/thrift"
	"time"
)

func AddFriend(ctx context.Context, req *thrift.AddFriendReq) (*thrift.AddFriendResp, error) {
	ok, err := mysql.CheckDuplicateRequest(req.UserIdSend, req.UserIdRecv)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &thrift.AddFriendResp{StatusCode: 1}, nil
	}

	userRelationRequest, err := mysql.CreateUserRelationRequest(req.UserIdSend, req.UserIdRecv)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	err = redis.BroadcastMessage(req.UserIdRecv, map[string]interface{}{
		"Id":         userRelationRequest.Id,
		"UserIdSend": req.UserIdSend,
		"UserIdRecv": req.UserIdRecv,
		"CreateTime": now.Unix(),
	}, redis.AddFriend)
	if err != nil {
		logrus.Warnf("redis AddFriend failed, req: %+v, err: %+v", req, err)
	}
	return &thrift.AddFriendResp{ID: userRelationRequest.Id, CreateTime: now.Unix()}, nil
}

func ReplyAddFriend(ctx context.Context, req *thrift.ReplyAddFriendReq) (*thrift.ReplyAddFriendResp, error) {
	var userIdSend, userIdRecv int64
	var err error
	now := time.Now()

	switch req.Status {
	case int32(entity.Accepted):
		userIdSend, userIdRecv, err = mysql.AcceptedAddFriend(req.ID)
	case int32(entity.Rejected):
		userIdSend, userIdRecv, err = mysql.RejectedAddFriend(req.ID)
	default:
		return nil, errors.New(fmt.Sprintf("req.Status invalid, Status: %d", req.Status))
	}
	if err != nil {
		return nil, err
	}
	id, err := mysql.CreateReplyAddFriend(req.ID, userIdRecv, userIdSend, now)
	if err != nil {
		return nil, err
	}
	//	mq message
	err = redis.BroadcastMessage(userIdSend, map[string]interface{}{
		"Id":         id,
		"ParentId":   req.ID,
		"UserIdSend": userIdRecv,
		"UserIdRecv": userIdSend,
		"CreateTime": now.Unix(),
		"Status":     req.Status,
	}, redis.ReplyAddFriend)
	if err != nil {
		logrus.Warnf("redis ReplyAddFriend failed, err: %+v", err)
	}
	return &thrift.ReplyAddFriendResp{ID: id, CreateTime: now.Unix()}, nil
}

func FetchFriendsList(ctx context.Context, req *thrift.FetchFriendsListReq) (*thrift.FetchFriendsListResp, error) {
	userIds, err := mysql.GetFriendsList(req.UserId)
	if err != nil {
		return nil, err
	}
	onlineUserIds, err := redis.GetUsers()
	if err != nil {
		return nil, err
	}
	onlineUserIdMap := make(map[int64]bool)
	for _, userId := range onlineUserIds {
		onlineUserIdMap[userId] = true
	}
	onlineFriendUserIds := make([]int64, 0)
	offlineFriendUserIds := make([]int64, 0)
	for _, userId := range userIds {
		if onlineUserIdMap[userId] {
			onlineFriendUserIds = append(onlineFriendUserIds, userId)
		} else {
			offlineFriendUserIds = append(offlineFriendUserIds, userId)
		}
	}
	return &thrift.FetchFriendsListResp{OnlineUserIds: onlineFriendUserIds, OfflineUserIds: offlineFriendUserIds}, nil
}
