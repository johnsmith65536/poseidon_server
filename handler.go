package main

import (
	"github.com/gin-gonic/gin"
	"poseidon/handler"
	"poseidon/middleware"
)

func initHttpServer(addr string) {
	r := gin.New()
	r.Use(gin.Logger(), middleware.Recover())

	r1 := r.Group("")
	r1.POST("/login", handler.Login)

	r2 := r.Group("", middleware.Auth)

	r2.GET("/heart_beat/:user_id", handler.HeartBeat)

	r2.POST("/message", handler.SendMessage)
	r2.GET("/message", handler.SyncMessage)
	r2.PUT("/message/status", handler.UpdateMessageStatus)
	r2.GET("/message/status", handler.FetchMessageStatus)

	r2.POST("/user", handler.CreateUser)
	r2.GET("/user/search", handler.SearchUser)

	r2.POST("/logout", handler.Logout)

	r2.GET("/friend/:user_id", handler.FetchFriendList)
	r2.POST("/friend", handler.AddFriend)
	r2.POST("/friend/reply", handler.ReplyAddFriend)
	r2.DELETE("/friend", handler.DeleteFriend)

	r2.GET("/sts_info/:user_id", handler.GetSTSInfo)

	r2.POST("/object", handler.CreateObject)

	r2.POST("/group", handler.CreateGroup)
	r2.GET("/group/search", handler.SearchGroup)
	r2.GET("/group/list/:user_id", handler.FetchGroupList)
	r2.GET("/group/last_read_msg_id/:user_id", handler.GetLastReadMsgId)
	r2.PUT("/group/last_read_msg_id/:user_id", handler.UpdateLastReadMsgId)

	r2.DELETE("/group/member", handler.DeleteMember)
	r2.DELETE("/group", handler.DeleteGroup)

	r2.GET("/group/member/:group_id", handler.FetchMemberList)
	r2.POST("/group/member/add", handler.AddGroup)
	r2.POST("/group/member/add/reply", handler.ReplyAddGroup)
	r2.POST("/group/member/invite", handler.InviteGroup)
	r2.POST("/group/member/invite/reply", handler.ReplyInviteGroup)

	r.GET("/ping", handler.Ping)

	r.Run(addr)
}
