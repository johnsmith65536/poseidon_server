package main

import (
	"github.com/gin-gonic/gin"
	"poseidon/handler"
)

func initHttpServer(addr string) {
	r := gin.Default()
	/*r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		body, err := ioutil.ReadAll(param.Request.Body)
		if err != nil {
			fmt.Println(err)
		}
		return fmt.Sprintf("%s %s",
			param.ClientIP,
			string(body),
		)
	}))*/
	r.GET("/heart_beat/:user_id", handler.HeartBeat)

	r.POST("/message", handler.SendMessage)
	r.GET("/message", handler.SyncMessage)
	r.PUT("/message/status", handler.UpdateMessageStatus)
	r.GET("/message/status", handler.FetchMessageStatus)

	r.POST("/user", handler.CreateUser)
	r.GET("/user/search", handler.SearchUser)

	r.POST("/login", handler.Login)
	r.POST("/logout", handler.Logout)

	r.GET("/friend/:user_id", handler.FetchFriendList)
	r.POST("/friend", handler.AddFriend)
	r.POST("/friend/reply", handler.ReplyAddFriend)
	r.DELETE("/friend", handler.DeleteFriend)

	r.GET("/sts_info/:user_id", handler.GetSTSInfo)

	r.POST("/object", handler.CreateObject)

	r.GET("/ping", handler.Ping)

	r.Run(addr)
}
