package handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/mysql"
	"strconv"
)

type CreateGroupReq struct {
	Owner   int64
	Name    string
	UserIds []int64
}

type CreateGroupResp struct {
	Id int64
	Status
}

/*
type SearchGroupReq struct {
	UserId int64
	Data   string
}*/

type Group struct {
	Id         int64
	Name       string
	CreateTime int64
	IsMember   bool
}

type SearchGroupResp struct {
	Groups []Group
	Status
}

type DeleteGroupReq struct {
	UserId  int64
	GroupId int64
}

type DeleteGroupResp struct {
	Status
}

func CreateGroup(c *gin.Context) {
	var err error
	var req CreateGroupReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	groupId, err := mysql.CreateGroup(req.Owner, req.Name)
	PanicIfError(err)
	err = mysql.AddGroupMember(groupId, append(req.UserIds, req.Owner))
	PanicIfError(err)
	c.JSON(200, CreateGroupResp{Id: groupId})
}

func SearchGroup(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	PanicIfError(err)

	data := c.Query("data")

	groupSearched, err := mysql.SearchGroup(data)
	PanicIfError(err)
	groupJoined, err := mysql.GetGroupList(userId)
	PanicIfError(err)
	groupMap := make(map[int64]bool)
	for _, groupId := range groupJoined {
		groupMap[groupId] = true
	}

	ret := make([]Group, 0)
	for _, group := range groupSearched {
		ret = append(ret, Group{Id: group.Id, Name: group.Name, CreateTime: group.CreateTime, IsMember: groupMap[group.Id]})
	}
	c.JSON(200, SearchGroupResp{Groups: ret})
}

func DeleteGroup(c *gin.Context) {
	var err error
	var req DeleteGroupReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	groupInfo, err := mysql.GetGroupInfo(req.GroupId)
	PanicIfError(err)
	if groupInfo.Owner != req.UserId {
		c.JSON(200, DeleteGroupResp{Status: Status{StatusCode: 1, StatusMessage: "only owner can operate"}})
		return
	}

	err = mysql.DeleteGroup(req.GroupId)
	PanicIfError(err)
	c.JSON(200, DeleteGroupResp{})
}
