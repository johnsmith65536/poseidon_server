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

type DeleteMemberReq struct {
	Operator int64
	GroupId  int64
	UserId   int64
}

type DeleteMemberResp struct {
	Status
}

type AddGroupReq struct {
	UserId  int64
	GroupId int64
}

type AddGroupResp struct {
	Id         int64
	CreateTime int64
	Status
}

type ReplyAddGroupReq struct {
	Id     int64
	Status int32
}

type ReplyAddGroupResp struct {
	Id         int64
	CreateTime int64
	Status
}

type InviteAddGroupReq struct {
	UserIdSend int64
	UserIdRecv int64
	GroupId    int64
}

type InviteAddGroupResp struct {
	Id         int64
	CreateTime int64
	Status
}

type ReplyInviteGroupReq struct {
	Id     int64
	Status int32
}

type ReplyInviteGroupResp struct {
	Id         int64
	CreateTime int64
	Status
}

/*type FetchGroupListReq struct {
	UserId int64
}*/

type FetchGroupListResp struct {
	Groups []*entity.Group
	Status
}

/*type FetchMemberListReq struct {
	GroupId int64
}*/

type FetchMemberListResp struct {
	OnlineUserIds  []int64
	OfflineUserIds []int64
	Status
}

/*type GetLastReadMsgIdReq struct {
	UserId int64
}*/

type GetLastReadMsgIdResp struct {
	LastReadMsgId map[int64]int64
	Status
}

type UpdateLastReadMsgIdReq struct {
	LastReadMsgId map[int64]int64
}

type UpdateLastReadMsgIdResp struct {
	Status
}

func AddGroup(c *gin.Context) {
	var err error
	var req AddGroupReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	ok, err := mysql.CheckDuplicateGroupRequest(req.UserId, req.GroupId)
	PanicIfError(err)
	if !ok {
		c.JSON(200, AddGroupResp{Status: Status{StatusCode: 1, StatusMessage: "duplicate request"}})
		return
	}

	isMember, err := mysql.CheckIsGroupMember(req.UserId, req.GroupId)
	PanicIfError(err)
	if isMember {
		c.JSON(200, AddGroupResp{Status: Status{StatusCode: 2, StatusMessage: "already group member"}})
		return
	}

	groupUserRequest, err := mysql.CreateGroupUserRequest(req.UserId, req.GroupId, entity.AddGroup)
	PanicIfError(err)

	group, err := mysql.GetGroupInfo(req.GroupId)
	PanicIfError(err)

	now := time.Now()
	err = redis.BroadcastMessage(group.Owner, map[string]interface{}{
		"Id":         groupUserRequest.Id,
		"UserIdSend": req.UserId,
		"GroupId":    req.GroupId,
		"CreateTime": now.Unix(),
	}, redis.AddGroup)
	if err != nil {
		logrus.Warnf("redis AddGroup failed, req: %+v, err: %+v", req, err)
	}
	c.JSON(200, AddGroupResp{Id: groupUserRequest.Id, CreateTime: now.Unix()})

}
func ReplyAddGroup(c *gin.Context) {
	var err error
	var req ReplyAddGroupReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)

	var userIdSend, groupId int64
	now := time.Now()

	switch req.Status {
	case int32(entity.Accepted):
		userIdSend, groupId, err = mysql.AcceptedAddGroup(req.Id)
	case int32(entity.Rejected):
		userIdSend, groupId, err = mysql.RejectedAddGroup(req.Id)
	default:
		PanicIfError(fmt.Errorf("req.Status invalid, Status: %d", req.Status))
		return
	}
	PanicIfError(err)

	group, err := mysql.GetGroupInfo(groupId)
	PanicIfError(err)

	id, err := mysql.CreateReplyAddGroup(req.Id, group.Owner, userIdSend, groupId, now, entity.AddGroup)
	PanicIfError(err)
	//	mq message
	err = redis.BroadcastMessage(userIdSend, map[string]interface{}{
		"Id":         id,
		"ParentId":   req.Id,
		"UserIdSend": group.Owner,
		"UserIdRecv": userIdSend,
		"GroupId":    groupId,
		"CreateTime": now.Unix(),
		"Status":     req.Status,
	}, redis.ReplyAddGroup)
	if err != nil {
		logrus.Warnf("redis ReplyAddGroup failed, err: %+v", err)
	}
	c.JSON(200, ReplyAddGroupResp{Id: id, CreateTime: now.Unix()})
}
func InviteGroup(c *gin.Context) {

}
func ReplyInviteGroup(c *gin.Context) {

}

func DeleteMember(c *gin.Context) {
	var err error
	var req DeleteMemberReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	groupInfo, err := mysql.GetGroupInfo(req.GroupId)
	PanicIfError(err)
	if groupInfo.Owner != req.Operator {
		c.JSON(200, DeleteGroupResp{Status: Status{StatusCode: 1, StatusMessage: "only owner can operate"}})
		return
	}
	err = mysql.DeleteGroupMember(req.GroupId, req.UserId)
	PanicIfError(err)
	c.JSON(200, DeleteMemberResp{})
}

func FetchGroupList(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	PanicIfError(err)
	groupIds, err := mysql.GetGroupList(userId)
	PanicIfError(err)
	groups, err := mysql.GetGroupInfos(groupIds)
	PanicIfError(err)
	c.JSON(200, FetchGroupListResp{Groups: groups})
}

func FetchMemberList(c *gin.Context) {
	var err error

	groupId, err := strconv.ParseInt(c.Param("group_id"), 10, 64)
	PanicIfError(err)

	userIds, err := mysql.GetMemberList(groupId)
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
	c.JSON(200, FetchMemberListResp{OnlineUserIds: onlineFriendUserIds, OfflineUserIds: offlineFriendUserIds})
}

func GetLastReadMsgId(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	PanicIfError(err)
	lastReadMsgId, err := mysql.GetLastReadMsgId(userId)
	PanicIfError(err)
	c.JSON(200, GetLastReadMsgIdResp{LastReadMsgId: lastReadMsgId})
}

func UpdateLastReadMsgId(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	PanicIfError(err)
	var req UpdateLastReadMsgIdReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	err = mysql.UpdateLastReadMsgId(userId, req.LastReadMsgId)
	PanicIfError(err)
	c.JSON(200, UpdateLastReadMsgIdResp{})
}
