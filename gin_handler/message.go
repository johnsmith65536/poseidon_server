package gin_handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/entity"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"strconv"
	"time"
)

type SendMessageReq struct {
	UserIdSend  int64
	IdRecv      int64
	Content     string
	ContentType int32
	MessageType int32
}

type SendMessageResp struct {
	Id         int64
	CreateTime int64
	Status
}

func SendMessage(c *gin.Context) {
	var req SendMessageReq
	var err error
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, SendMessageResp{Status: Status{StatusCode: 1, StatusMessage: err.Error()}})
		return
	}
	msg, err := mysql.WriteMessage(req.UserIdSend, req.IdRecv, 0, req.Content, time.Now(), req.ContentType, req.MessageType, false)
	if err != nil {
		c.JSON(200, SendMessageResp{Status: Status{StatusCode: 1, StatusMessage: err.Error()}})
		return
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
			c.JSON(200, SendMessageResp{Status: Status{StatusCode: 1, StatusMessage: err.Error()}})
			return
		}
		object, err := mysql.GetObject(objId)
		if err != nil {
			c.JSON(200, SendMessageResp{Status: Status{StatusCode: 1, StatusMessage: err.Error()}})
			return
		}
		broadcastMsg["ObjectETag"] = object.ETag
		broadcastMsg["ObjectName"] = object.Name
	}
	go redis.BroadcastMessage(req.IdRecv, broadcastMsg, redis.Chat)
	c.JSON(200, SendMessageResp{Id: msg.Id, CreateTime: msg.CreateTime})
}
