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

func SearchUser(ctx context.Context, req *thrift.SearchUserReq) (*thrift.SearchUserResp, error) {
	var users []*thrift.User
	var userIds []int64
	userInfos, err := mysql.SearchUser(req.Data)
	if err != nil {
		return nil, err
	}
	for _, userInfo := range userInfos {
		userIds = append(userIds, userInfo.Id)
	}
	friendUserIds, err := mysql.GetRelation(req.UserId, userIds)
	if err != nil {
		return nil, err
	}
	friendMap := make(map[int64]bool)
	for _, friendUserId := range friendUserIds {
		friendMap[friendUserId] = true
	}
	for _, userInfo := range userInfos {
		users = append(users, &thrift.User{ID: userInfo.Id, NickName: userInfo.NickName, LastOnlineTime: userInfo.LastOnlineTime.Unix(), IsFriend: friendMap[userInfo.Id]})
	}
	return &thrift.SearchUserResp{Users: users}, nil
}
