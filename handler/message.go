package handler

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"poseidon/entity"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"poseidon/utils"
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
	PanicIfError(err)


	isFriend,err := mysql.CheckIsFriend(req.UserIdSend, req.IdRecv)
	PanicIfError(err)
	if !isFriend {
		c.JSON(200, SendMessageResp{Status: Status{StatusCode: 1, StatusMessage: "is not friend"}})
		return
	}

	byteContent, err := base64.StdEncoding.DecodeString(req.Content)
	PanicIfError(err)
	rawContent, err := utils.UnGzip(byteContent)
	PanicIfError(err)
	msg, err := mysql.WriteMessage(req.UserIdSend, req.IdRecv, 0, byteContent, time.Now(), req.ContentType, req.MessageType, false)
	PanicIfError(err)
	broadcastMsg := map[string]interface{}{
		"Id":          msg.Id,
		"UserIdSend":  msg.UserIdSend,
		"UserIdRecv":  msg.UserIdRecv,
		"GroupId":     msg.GroupId,
		"Content":     string(rawContent),
		"CreateTime":  msg.CreateTime,
		"ContentType": msg.ContentType,
		"MsgType":     msg.MsgType,
	}
	switch req.ContentType {
	case int32(entity.Text):
		;
	case int32(entity.ObjectData):
		objId, err := strconv.ParseInt(string(rawContent), 10, 64)
		PanicIfError(err)
		object, err := mysql.GetObject(objId)
		PanicIfError(err)
		broadcastMsg["ObjectETag"] = object.ETag
		broadcastMsg["ObjectName"] = object.Name
	case int32(entity.Vibration):
		;
	case int32(entity.Image):
		;
	}
	go redis.BroadcastMessage(req.IdRecv, broadcastMsg, redis.Chat)
	c.JSON(200, SendMessageResp{Id: msg.Id, CreateTime: msg.CreateTime})
}

func SyncMessage(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	PanicIfError(err)

	messageId, err := strconv.ParseInt(c.Query("message_id"), 10, 64)
	PanicIfError(err)

	userRelationId, err := strconv.ParseInt(c.Query("user_relation_id"), 10, 64)
	PanicIfError(err)

	messages, err := mysql.SyncMessage(userId, messageId)
	PanicIfError(err)
	userRelations, err := mysql.SyncUserRelationRequest(userId, userRelationId)
	PanicIfError(err)
	var objIds []int64
	for _, message := range messages {
		if message.ContentType == int32(entity.ObjectData) {
			content, err := utils.UnGzip(message.Content)
			PanicIfError(err)
			objId, err := strconv.ParseInt(string(content), 10, 64)
			PanicIfError(err)
			objIds = append(objIds, objId)
		}
	}
	objects, err := mysql.SyncObject(objIds)
	PanicIfError(err)
	lastOnlineTime, err := mysql.GetLastOnlineTime(userId)
	PanicIfError(err)
	c.JSON(200, SyncMessageResp{Messages: messages, UserRelations: userRelations, Objects: objects, LastOnlineTime: lastOnlineTime})
}

func UpdateMessageStatus(c *gin.Context) {
	var req UpdateMessageStatusReq
	var err error
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	err = mysql.UpdateMessageStatus(req.MessageIds, req.UserRelationRequestIds)
	PanicIfError(err)
	c.JSON(200, UpdateMessageStatusResp{})
}

func FetchMessageStatus(c *gin.Context) {
	var err error
	var messageIds []int64
	var userRelationRequestIds []int64
	messageIdsStr := c.QueryArray("message_ids")
	for _, messageIdStr := range messageIdsStr {
		messageId, err := strconv.ParseInt(messageIdStr, 10, 64)
		PanicIfError(err)
		messageIds = append(messageIds, messageId)
	}

	userRelationRequestIdsStr := c.QueryArray("user_relation_request_ids")
	for _, userRelationRequestIdStr := range userRelationRequestIdsStr {
		userRelationRequestId, err := strconv.ParseInt(userRelationRequestIdStr, 10, 64)
		PanicIfError(err)
		userRelationRequestIds = append(userRelationRequestIds, userRelationRequestId)
	}

	messageStatus, err := mysql.GetMessageStatus(messageIds)
	PanicIfError(err)
	relationStatus, err := mysql.GetRelationStatus(userRelationRequestIds)
	PanicIfError(err)
	c.JSON(200, FetchMessageStatusResp{MessageIds: messageStatus, UserRelationRequestIds: relationStatus})
}
