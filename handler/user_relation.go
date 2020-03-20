package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"poseidon/entity"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"strconv"
	"time"
)

type AddFriendReq struct {
	UserIdSend int64
	UserIdRecv int64
}

type AddFriendResp struct {
	Id         int64
	CreateTime int64
	//StatusCode int32 `thrift:"StatusCode,255" db:"StatusCode" json:"StatusCode"`
	Status
}

type ReplyAddFriendReq struct {
	Id     int64
	Status int32
}

type ReplyAddFriendResp struct {
	Id         int64
	CreateTime int64
	Status
}

/*type FetchFriendListReq struct {
	UserId int64 `thrift:"UserId,1" db:"UserId" json:"UserId"`
}*/

type FetchFriendListResp struct {
	OnlineUserIds  []int64
	OfflineUserIds []int64
	Status
}

/*type DeleteFriendReq struct {
	UserIdSend int64 `thrift:"UserIdSend,1" db:"UserIdSend" json:"UserIdSend"`
	UserIdRecv int64 `thrift:"UserIdRecv,2" db:"UserIdRecv" json:"UserIdRecv"`
}
*/
type DeleteFriendResp struct {
	Status
}

func FetchFriendList(c *gin.Context) {
	var err error

	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	PanicIfError(err)

	userIds, err := mysql.GetFriendsList(userId)
	PanicIfError(err)
	onlineUserIds, err := redis.GetUsers()
	PanicIfError(err)
	onlineUserIdMap := make(map[int64]bool)
	for _, userId := range onlineUserIds {
		onlineUserIdMap[userId] = true
	}
	onlineFriendUserIds := make([]int64, 0)
	offlineFriendUserIds := make([]int64, 0)
	for _, userId := range userIds {
		if onlineUserIdMap[userId] {
			onlineFriendUserIds = append(onlineFriendUserIds, userId)
		} else {
			offlineFriendUserIds = append(offlineFriendUserIds, userId)
		}
	}
	c.JSON(200, FetchFriendListResp{OnlineUserIds: onlineFriendUserIds, OfflineUserIds: offlineFriendUserIds})
}

func AddFriend(c *gin.Context) {
	var err error
	var req AddFriendReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	ok, err := mysql.CheckDuplicateRequest(req.UserIdSend, req.UserIdRecv)
	PanicIfError(err)
	if !ok {
		c.JSON(200, AddFriendResp{Status: Status{StatusCode: 1, StatusMessage: "duplicate request"}})
		return
	}

	isFriend, err := mysql.CheckIsFriend(req.UserIdSend, req.UserIdRecv)
	PanicIfError(err)
	if isFriend {
		c.JSON(200, AddFriendResp{Status: Status{StatusCode: 2, StatusMessage: "already friend"}})
		return
	}

	userRelationRequest, err := mysql.CreateUserRelationRequest(req.UserIdSend, req.UserIdRecv)
	PanicIfError(err)
	now := time.Now()
	err = redis.BroadcastMessage(req.UserIdRecv, map[string]interface{}{
		"Id":         userRelationRequest.Id,
		"UserIdSend": req.UserIdSend,
		"UserIdRecv": req.UserIdRecv,
		"CreateTime": now.Unix(),
	}, redis.AddFriend)
	if err != nil {
		logrus.Warnf("redis AddFriend failed, req: %+v, err: %+v", req, err)
	}
	c.JSON(200, AddFriendResp{Id: userRelationRequest.Id, CreateTime: now.Unix()})
}

func ReplyAddFriend(c *gin.Context) {
	var err error
	var req ReplyAddFriendReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)

	var userIdSend, userIdRecv int64
	now := time.Now()

	switch req.Status {
	case int32(entity.Accepted):
		userIdSend, userIdRecv, err = mysql.AcceptedAddFriend(req.Id)
	case int32(entity.Rejected):
		userIdSend, userIdRecv, err = mysql.RejectedAddFriend(req.Id)
	default:
		PanicIfError(fmt.Errorf("req.Status invalid, Status: %d", req.Status))
		return
	}
	PanicIfError(err)
	id, err := mysql.CreateReplyAddFriend(req.Id, userIdRecv, userIdSend, now)
	PanicIfError(err)
	//	mq message
	err = redis.BroadcastMessage(userIdSend, map[string]interface{}{
		"Id":         id,
		"ParentId":   req.Id,
		"UserIdSend": userIdRecv,
		"UserIdRecv": userIdSend,
		"CreateTime": now.Unix(),
		"Status":     req.Status,
	}, redis.ReplyAddFriend)
	if err != nil {
		logrus.Warnf("redis ReplyAddFriend failed, err: %+v", err)
	}
	c.JSON(200, ReplyAddFriendResp{Id: id, CreateTime: now.Unix()})
}

func DeleteFriend(c *gin.Context) {
	var err error
	userIdSend, err := strconv.ParseInt(c.Query("user_id_send"), 10, 64)
	PanicIfError(err)

	userIdRecv, err := strconv.ParseInt(c.Query("user_id_recv"), 10, 64)
	PanicIfError(err)

	err = mysql.DeleteFriend(userIdSend, userIdRecv)
	PanicIfError(err)
	c.JSON(200, DeleteFriendResp{})
}
