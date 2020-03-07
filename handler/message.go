package handler

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

/*type SyncMessageReq struct {
	UserId         int64
	MessageId      int64
	UserRelationId int64
}*/

type SyncMessageResp struct {
	Messages       []*entity.Message
	UserRelations  []*entity.UserRelationRequest
	Objects        []*entity.Object
	LastOnlineTime int64
	Status
}

type UpdateMessageStatusReq struct {
	MessageIds             map[int64]int32
	UserRelationRequestIds map[int64]int32
}

type UpdateMessageStatusResp struct {
	Status
}

/*type FetchMessageStatusReq struct {
	MessageIds []int64 `thrift:"MessageIds,1" db:"MessageIds" json:"MessageIds"`
	UserRelationRequestIds []int64 `thrift:"UserRelationRequestIds,2" db:"UserRelationRequestIds" json:"UserRelationRequestIds"`
}*/

type FetchMessageStatusResp struct {
	MessageIds             map[int64]int32
	UserRelationRequestIds map[int64]int32
	Status
}

func SendMessage(c *gin.Context) {
	var req SendMessageReq
	var err error
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, SendMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	msg, err := mysql.WriteMessage(req.UserIdSend, req.IdRecv, 0, req.Content, time.Now(), req.ContentType, req.MessageType, false)
	if err != nil {
		c.JSON(200, SendMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
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
			c.JSON(200, SendMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
			return
		}
		object, err := mysql.GetObject(objId)
		if err != nil {
			c.JSON(200, SendMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
			return
		}
		broadcastMsg["ObjectETag"] = object.ETag
		broadcastMsg["ObjectName"] = object.Name
	}
	go redis.BroadcastMessage(req.IdRecv, broadcastMsg, redis.Chat)
	c.JSON(200, SendMessageResp{Id: msg.Id, CreateTime: msg.CreateTime})
}

func SyncMessage(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(200, SyncMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}

	messageId, err := strconv.ParseInt(c.Query("message_id"), 10, 64)
	if err != nil {
		c.JSON(200, SyncMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}

	userRelationId, err := strconv.ParseInt(c.Query("user_relation_id"), 10, 64)
	if err != nil {
		c.JSON(200, SyncMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}

	messages, err := mysql.SyncMessage(userId, messageId)
	if err != nil {
		c.JSON(200, SyncMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	userRelations, err := mysql.SyncUserRelationRequest(userId, userRelationId)
	if err != nil {
		c.JSON(200, SyncMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	var objIds []int64
	for _, message := range messages {
		if message.ContentType == int32(entity.ObjectData) {
			objId, err := strconv.ParseInt(message.Content, 10, 64)
			if err != nil {
				c.JSON(200, SyncMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
				return
			}
			objIds = append(objIds, objId)
		}
	}
	objects, err := mysql.SyncObject(objIds)
	if err != nil {
		c.JSON(200, SyncMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	lastOnlineTime, err := mysql.GetLastOnlineTime(userId)
	if err != nil {
		c.JSON(200, SyncMessageResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	c.JSON(200, SyncMessageResp{Messages: messages, UserRelations: userRelations, Objects: objects, LastOnlineTime: lastOnlineTime})
}

func UpdateMessageStatus(c *gin.Context) {
	var req UpdateMessageStatusReq
	var err error
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, UpdateMessageStatusResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	err = mysql.UpdateMessageStatus(req.MessageIds, req.UserRelationRequestIds)
	if err != nil {
		c.JSON(200, UpdateMessageStatusResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	c.JSON(200, UpdateMessageStatusResp{})
}

func FetchMessageStatus(c *gin.Context) {
	var err error
	var messageIds []int64
	var userRelationRequestIds []int64
	messageIdsStr := c.QueryArray("message_ids")
	if err != nil {
		c.JSON(200, FetchMessageStatusResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	for _, messageIdStr := range messageIdsStr {
		messageId, err := strconv.ParseInt(messageIdStr, 10, 64)
		if err != nil {
			c.JSON(200, FetchMessageStatusResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
			return
		}
		messageIds = append(messageIds, messageId)
	}

	userRelationRequestIdsStr := c.QueryArray("user_relation_request_ids")
	if err != nil {
		c.JSON(200, FetchMessageStatusResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	for _, userRelationRequestIdStr := range userRelationRequestIdsStr {
		userRelationRequestId, err := strconv.ParseInt(userRelationRequestIdStr, 10, 64)
		if err != nil {
			c.JSON(200, FetchMessageStatusResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
			return
		}
		userRelationRequestIds = append(userRelationRequestIds, userRelationRequestId)
	}

	messageStatus, err := mysql.GetMessageStatus(messageIds)
	if err != nil {
		c.JSON(200, FetchMessageStatusResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	relationStatus, err := mysql.GetRelationStatus(userRelationRequestIds)
	if err != nil {
		c.JSON(200, FetchMessageStatusResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	c.JSON(200, FetchMessageStatusResp{MessageIds:messageStatus,UserRelationRequestIds:relationStatus})
}
