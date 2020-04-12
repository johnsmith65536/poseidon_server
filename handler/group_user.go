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
	GroupId int64
	UserId  int64
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

type InviteGroupReq struct {
	GroupId int64
	UserIds  []int64
}

type InviteGroupResp struct {
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
	OnlineUsers  []*entity.User
	OfflineUsers []*entity.User
	Status
}

/*type GetGroupLastReadMsgIdReq struct {
	UserId int64
}*/

type GetGroupLastReadMsgIdResp struct {
	LastReadMsgId map[int64]int64
	Status
}

type UpdateGroupLastReadMsgIdReq struct {
	LastReadMsgId map[int64]int64
}

type UpdateGroupLastReadMsgIdResp struct {
	Status
}

/*type InviteGroupFriendListReq struct {
	UserId  int64
	GroupId int64
}*/

type InviteGroupFriendListResp struct {
	NotInGroupUsers []*entity.User
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
	var err error
	var req InviteGroupReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	err = mysql.InviteGroup(req.GroupId, req.UserIds)
	PanicIfError(err)
	c.JSON(200, InviteGroupResp{})
}

func DeleteMember(c *gin.Context) {
	var err error
	var req DeleteMemberReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
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

	var res map[int64]string
	onlineUsers, offlineUsers := make([]*entity.User, 0), make([]*entity.User, 0)
	res, err = mysql.GetUserNickNames(onlineFriendUserIds)
	PanicIfError(err)
	for _, userId := range onlineFriendUserIds {
		onlineUsers = append(onlineUsers, &entity.User{Id: userId, NickName: res[userId]})
	}

	res, err = mysql.GetUserNickNames(offlineFriendUserIds)
	PanicIfError(err)
	for _, userId := range offlineFriendUserIds {
		offlineUsers = append(offlineUsers, &entity.User{Id: userId, NickName: res[userId]})
	}

	c.JSON(200, FetchMemberListResp{OnlineUsers: onlineUsers, OfflineUsers: offlineUsers})
}

func GetGroupLastReadMsgId(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	PanicIfError(err)
	lastReadMsgId, err := mysql.GetGroupLastReadMsgId(userId)
	PanicIfError(err)
	c.JSON(200, GetGroupLastReadMsgIdResp{LastReadMsgId: lastReadMsgId})
}

func UpdateGroupLastReadMsgId(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	PanicIfError(err)
	var req UpdateGroupLastReadMsgIdReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	err = mysql.UpdateGroupLastReadMsgId(userId, req.LastReadMsgId)
	PanicIfError(err)
	c.JSON(200, UpdateGroupLastReadMsgIdResp{})
}

func InviteGroupFriendList(c *gin.Context) {
	var err error
	notInGroupUsers := make([]*entity.User, 0)
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	PanicIfError(err)
	groupId, err := strconv.ParseInt(c.Query("group_id"), 10, 64)
	PanicIfError(err)
	userIds, err := mysql.GetFriendsList(userId)
	PanicIfError(err)

	nickName, err := mysql.GetUserNickNames(userIds)
	PanicIfError(err)
	for _, friendId := range userIds {
		ret, err := mysql.CheckIsGroupMember(friendId, groupId)
		PanicIfError(err)
		if !ret {
			notInGroupUsers = append(notInGroupUsers, &entity.User{Id: friendId, NickName: nickName[friendId]})
		}
	}
	c.JSON(200, InviteGroupFriendListResp{NotInGroupUsers: notInGroupUsers})
}
