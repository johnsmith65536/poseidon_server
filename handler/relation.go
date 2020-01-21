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
	userRelationRequest, err := mysql.CreateUserRelationRequest(req.UserIdSend, req.UserIdRecv)
	if err != nil {
		return nil, err
	}
	err = redis.BroadcastMessage(req.UserIdRecv, map[string]interface{}{
		"Id":         userRelationRequest.Id,
		"UserIdSend": req.UserIdSend,
		"UserIdRecv": req.UserIdRecv,
		"CreateTime": time.Now().Unix(),
	}, redis.AddFriend)
	if err != nil {
		logrus.Warnf("redis AddFriend failed, req: %+v, err: %+v", req, err)
	}
	return &thrift.AddFriendResp{}, nil
}

func ReplyAddFriend(ctx context.Context, req *thrift.ReplyAddFriendReq) (*thrift.ReplyAddFriendResp, error) {
	var err error
	switch req.Status {
	case int32(entity.Accepted):
		err = mysql.AcceptedAddFriend(req.ID)
	case int32(entity.Rejected):
		err = mysql.RejectedAddFriend(req.ID)
	default:
		return nil, errors.New(fmt.Sprintf("req.Status invalid, Status: %d", req.Status))
	}
	if err != nil {
		return nil, err
	}
	return &thrift.ReplyAddFriendResp{}, nil
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
