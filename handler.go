package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"poseidon/gin_handler"
	"poseidon/handler"
	"poseidon/thrift"
)

type ServiceImpl struct{}

func (fdi *ServiceImpl) HeartBeat(ctx context.Context, req *thrift.HeartBeatReq) (resp *thrift.HeartBeatResp, err error) {
	resp, err = handler.HeartBeat(ctx, req)
	if err != nil {
		logrus.Errorf("HeartBeat failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) CreateUser(ctx context.Context, req *thrift.CreateUserReq) (resp *thrift.CreateUserResp, err error) {
	resp, err = handler.CreateUser(ctx, req)
	if err != nil {
		logrus.Errorf("CreateUser failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) Login(ctx context.Context, req *thrift.LoginReq) (resp *thrift.LoginResp, err error) {
	resp, err = handler.Login(ctx, req)
	if err != nil {
		logrus.Errorf("Login failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) Logout(ctx context.Context, req *thrift.LogoutReq) (resp *thrift.LogoutResp, err error) {
	resp, err = handler.Logout(ctx, req)
	if err != nil {
		logrus.Errorf("Logout failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) SendMessage(ctx context.Context, req *thrift.SendMessageReq) (resp *thrift.SendMessageResp, err error) {
	resp, err = handler.SendMessage(ctx, req)
	if err != nil {
		logrus.Errorf("SendMessage failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}
func (fdi *ServiceImpl) FetchFriendsList(ctx context.Context, req *thrift.FetchFriendsListReq) (resp *thrift.FetchFriendsListResp, err error) {
	resp, err = handler.FetchFriendsList(ctx, req)
	if err != nil {
		logrus.Errorf("FetchFriendsList failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) SyncMessage(ctx context.Context, req *thrift.SyncMessageReq) (resp *thrift.SyncMessageResp, err error) {
	resp, err = handler.SyncMessage(ctx, req)
	if err != nil {
		logrus.Errorf("SyncMessage failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}
func (fdi *ServiceImpl) CreateObject(ctx context.Context, req *thrift.CreateObjectReq) (resp *thrift.CreateObjectResp, err error) {
	resp, err = handler.CreateObject(ctx, req)
	if err != nil {
		logrus.Errorf("CreateObject failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}
func (fdi *ServiceImpl) AddFriend(ctx context.Context, req *thrift.AddFriendReq) (resp *thrift.AddFriendResp, err error) {
	resp, err = handler.AddFriend(ctx, req)
	if err != nil {
		logrus.Errorf("AddFriend failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}
func (fdi *ServiceImpl) ReplyAddFriend(ctx context.Context, req *thrift.ReplyAddFriendReq) (resp *thrift.ReplyAddFriendResp, err error) {
	resp, err = handler.ReplyAddFriend(ctx, req)
	if err != nil {
		logrus.Errorf("ReplyAddFriend failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) UpdateMessageStatus(ctx context.Context, req *thrift.UpdateMessageStatusReq) (resp *thrift.UpdateMessageStatusResp, err error) {
	resp, err = handler.UpdateMessageStatus(ctx, req)
	if err != nil {
		logrus.Errorf("UpdateMessageStatus failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) SearchUser(ctx context.Context, req *thrift.SearchUserReq) (resp *thrift.SearchUserResp, err error) {
	resp, err = handler.SearchUser(ctx, req)
	if err != nil {
		logrus.Errorf("SearchUser failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) DeleteFriend(ctx context.Context, req *thrift.DeleteFriendReq) (resp *thrift.DeleteFriendResp, err error) {
	resp, err = handler.DeleteFriend(ctx, req)
	if err != nil {
		logrus.Errorf("DeleteFriend failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) GetSTSInfo(ctx context.Context, req *thrift.GetSTSInfoReq) (resp *thrift.GetSTSInfoResp, err error) {
	resp, err = handler.GetSTSInfo(ctx, req)
	if err != nil {
		logrus.Errorf("GetSTSInfo failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func (fdi *ServiceImpl) FetchMessageStatus(ctx context.Context, req *thrift.FetchMessageStatusReq) (resp *thrift.FetchMessageStatusResp, err error) {
	resp, err = handler.FetchMessageStatus(ctx, req)
	if err != nil {
		logrus.Errorf("FetchMessageStatus failed, err: %+s", err)
	}
	logrus.Info(req)
	logrus.Info(resp)
	return resp, err
}

func initHttpServer(addr string) {
	r := gin.Default()
	r.GET("/heart_beat/:user_id", gin_handler.HeartBeat)
	r.POST("/message", gin_handler.SendMessage)
	r.Run(addr)
}
