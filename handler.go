package main

import (
	"context"
	"fmt"
)

type ServiceImpl struct{}

func (fdi *ServiceImpl) HeartBeat(ctx context.Context, req *HeartBeatReq) (resp *HeartBeatResp, err error) {
	fmt.Printf("%+v\n", req)
	return &HeartBeatResp{}, nil
}

func (fdi *ServiceImpl) CreateUser(ctx context.Context, req *CreateUserReq) (resp *CreateUserResp, err error) {
	fmt.Printf("%+v\n", req)
	return &CreateUserResp{}, nil
}

func (fdi *ServiceImpl) Login(ctx context.Context, req *LoginReq) (resp *LoginResp, err error) {
	fmt.Printf("%+v\n", req)
	return &LoginResp{}, nil
}

func (fdi *ServiceImpl) Logout(ctx context.Context, req *LogoutReq) (resp *LogoutResp, err error) {
	fmt.Printf("%+v\n", req)
	return &LogoutResp{}, nil
}

func (fdi *ServiceImpl) SendMessage(ctx context.Context, req *SendMessageReq) (resp *SendMessageResp, err error) {
	fmt.Printf("%+v\n", req)
	return &SendMessageResp{}, nil
}
func (fdi *ServiceImpl) FetchFriendsList(ctx context.Context, req *FetchFriendsListReq) (resp *FetchFriendsListResp, err error) {
	fmt.Printf("%+v\n", req)
	return &FetchFriendsListResp{}, nil
}

func (fdi *ServiceImpl) FetchOfflineMessage(ctx context.Context, req *FetchOfflineMessageReq) (resp *FetchOfflineMessageResp, err error) {
	fmt.Printf("%+v\n", req)
	return &FetchOfflineMessageResp{}, nil
}
func (fdi *ServiceImpl) UploadObject(ctx context.Context, req *UploadObjectReq) (resp *UploadObjectResp, err error) {
	fmt.Printf("%+v\n", req)
	return &UploadObjectResp{}, nil
}
func (fdi *ServiceImpl) AddFriend(ctx context.Context, req *AddFriendReq) (resp *AddFriendResp, err error) {
	fmt.Printf("%+v\n", req)
	return &AddFriendResp{}, nil
}
func (fdi *ServiceImpl) ReplyAddFriend(ctx context.Context, req *ReplyAddFriendReq) (resp *ReplyAddFriendResp, err error) {
	fmt.Printf("%+v\n", req)
	return &ReplyAddFriendResp{}, nil
}




