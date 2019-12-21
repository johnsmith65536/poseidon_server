package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"poseidon/handler"
	"poseidon/thrift"
)

type ServiceImpl struct{}

func (fdi *ServiceImpl) HeartBeat(ctx context.Context, req *thrift.HeartBeatReq) (resp *thrift.HeartBeatResp, err error) {
	resp, err = handler.HeartBeat(ctx, req)
	if err != nil {
		logrus.Errorf("HeartBeat failed, err: %+s", err)
	}
	return resp, err
}

func (fdi *ServiceImpl) CreateUser(ctx context.Context, req *thrift.CreateUserReq) (resp *thrift.CreateUserResp, err error) {
	resp, err = handler.CreateUser(ctx, req)
	if err != nil {
		logrus.Errorf("CreateUser failed, err: %+s", err)
	}
	return resp, err
}

func (fdi *ServiceImpl) Login(ctx context.Context, req *thrift.LoginReq) (resp *thrift.LoginResp, err error) {
	resp, err = handler.Login(ctx, req)
	if err != nil {
		logrus.Errorf("Login failed, err: %+s", err)
	}
	return resp, err
}

func (fdi *ServiceImpl) Logout(ctx context.Context, req *thrift.LogoutReq) (resp *thrift.LogoutResp, err error) {
	resp, err = handler.Logout(ctx, req)
	if err != nil {
		logrus.Errorf("Logout failed, err: %+s", err)
	}
	return resp, err
}

func (fdi *ServiceImpl) SendMessage(ctx context.Context, req *thrift.SendMessageReq) (resp *thrift.SendMessageResp, err error) {
	fmt.Printf("%+v\n", req)
	return &thrift.SendMessageResp{}, nil
}
func (fdi *ServiceImpl) FetchFriendsList(ctx context.Context, req *thrift.FetchFriendsListReq) (resp *thrift.FetchFriendsListResp, err error) {
	fmt.Printf("%+v\n", req)
	return &thrift.FetchFriendsListResp{}, nil
}

func (fdi *ServiceImpl) FetchOfflineMessage(ctx context.Context, req *thrift.FetchOfflineMessageReq) (resp *thrift.FetchOfflineMessageResp, err error) {
	fmt.Printf("%+v\n", req)
	return &thrift.FetchOfflineMessageResp{}, nil
}
func (fdi *ServiceImpl) UploadObject(ctx context.Context, req *thrift.UploadObjectReq) (resp *thrift.UploadObjectResp, err error) {
	fmt.Printf("%+v\n", req)
	return &thrift.UploadObjectResp{}, nil
}
func (fdi *ServiceImpl) AddFriend(ctx context.Context, req *thrift.AddFriendReq) (resp *thrift.AddFriendResp, err error) {
	fmt.Printf("%+v\n", req)
	return &thrift.AddFriendResp{}, nil
}
func (fdi *ServiceImpl) ReplyAddFriend(ctx context.Context, req *thrift.ReplyAddFriendReq) (resp *thrift.ReplyAddFriendResp, err error) {
	fmt.Printf("%+v\n", req)
	return &thrift.ReplyAddFriendResp{}, nil
}
