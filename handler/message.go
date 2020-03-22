package handler

import (
	"encoding/base64"
	"errors"
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
	GroupUserId int64
}*/

type SyncMessageResp struct {
	Messages       []*entity.Message
	UserRelations  []*entity.UserRelationRequest
	Objects        []*entity.Object
	GroupUsers     []*entity.GroupUserRequest
	LastOnlineTime int64
	Status
}

func SendMessage(c *gin.Context) {

	var req SendMessageReq
	var err error
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)

	byteContent, err := base64.StdEncoding.DecodeString(req.Content)
	PanicIfError(err)
	rawContent, err := utils.UnGzip(byteContent)
	PanicIfError(err)
	var msg *entity.Message

	switch req.MessageType {
	case int32(entity.PrivateChat):
		isFriend, err := mysql.CheckIsFriend(req.UserIdSend, req.IdRecv)
		PanicIfError(err)
		if !isFriend {
			c.JSON(200, SendMessageResp{Status: Status{StatusCode: 1, StatusMessage: "is not friend"}})
			return
		}
		msg, err = mysql.WriteMessage(req.UserIdSend, req.IdRecv, 0, byteContent, time.Now(), req.ContentType, req.MessageType)
		PanicIfError(err)
		err = mysql.UpdateFriendLastReadMsgId(req.UserIdSend, map[int64]int64{req.IdRecv: msg.Id})
		PanicIfError(err)
	case int32(entity.GroupChat):
		isGroupMember, err := mysql.CheckIsGroupMember(req.UserIdSend, req.IdRecv)
		PanicIfError(err)
		if !isGroupMember {
			c.JSON(200, SendMessageResp{Status: Status{StatusCode: 1, StatusMessage: "is not group member"}})
			return
		}
		msg, err = mysql.WriteMessage(req.UserIdSend, 0, req.IdRecv, byteContent, time.Now(), req.ContentType, req.MessageType)
		PanicIfError(err)
		err = mysql.UpdateGroupLastReadMsgId(req.UserIdSend, map[int64]int64{req.IdRecv: msg.Id})
		PanicIfError(err)
	default:
		PanicIfError(errors.New("unknown msgType"))
	}



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

	switch req.MessageType {
	case int32(entity.PrivateChat):
		err = redis.BroadcastMessage(req.IdRecv, broadcastMsg, redis.Chat)
		PanicIfError(err)
	case int32(entity.GroupChat):
		members, err := mysql.GetMemberList(req.IdRecv)
		PanicIfError(err)
		onlineUsers, err := redis.GetUsers()
		PanicIfError(err)
		userIds := utils.Intersection(members, onlineUsers)
		for userId := range userIds {
			if userId == req.UserIdSend {
				continue
			}
			err = redis.BroadcastMessage(userId, broadcastMsg, redis.Chat)
			PanicIfError(err)
		}
	default:
		PanicIfError(errors.New("unknown msgType"))
	}

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

	groupUserId, err := strconv.ParseInt(c.Query("group_user_id"), 10, 64)
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
	groupUsers, err := mysql.SyncGroupUserRequest(userId, groupUserId)
	PanicIfError(err)
	lastOnlineTime, err := mysql.GetLastOnlineTime(userId)
	PanicIfError(err)
	c.JSON(200, SyncMessageResp{Messages: messages, UserRelations: userRelations, Objects: objects, GroupUsers: groupUsers, LastOnlineTime: lastOnlineTime})
}

