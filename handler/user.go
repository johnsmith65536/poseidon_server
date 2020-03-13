package handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/mysql"
	"poseidon/utils"
	"strconv"
	"time"
)

type User struct {
	Id             int64
	NickName       string
	LastOnlineTime int64
	IsFriend       bool
}

type CreateUserReq struct {
	Password string
	NickName string
}

type CreateUserResp struct {
	UserId int64
	Status
}

/*type SearchUserReq struct {
	UserId int64 `thrift:"UserId,1" db:"UserId" json:"UserId"`
	Data string `thrift:"Data,2" db:"Data" json:"Data"`
}*/

type SearchUserResp struct {
	Users []*User
	Status
}

func CreateUser(c *gin.Context) {
	var err error
	var req CreateUserReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, CreateUserResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}

	userId := utils.GenerateId(4)
	err = mysql.CreateUser(userId, req.Password, req.NickName, time.Now())
	if err != nil {
		c.JSON(200, CreateUserResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	c.JSON(200, CreateUserResp{UserId: userId})
}

func SearchUser(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(200, SearchUserResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}

	data := c.Query("data")

	users := make([]*User, 0)
	userIds := make([]int64, 0)
	userInfos, err := mysql.SearchUser(data)
	if err != nil {
		c.JSON(200, SearchUserResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	for _, userInfo := range userInfos {
		userIds = append(userIds, userInfo.Id)
	}
	friendUserIds, err := mysql.GetRelation(userId, userIds)
	if err != nil {
		c.JSON(200, SearchUserResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	friendMap := make(map[int64]bool)
	for _, friendUserId := range friendUserIds {
		friendMap[friendUserId] = true
	}
	for _, userInfo := range userInfos {
		users = append(users, &User{Id: userInfo.Id, NickName: userInfo.NickName, LastOnlineTime: userInfo.LastOnlineTime.Unix(), IsFriend: friendMap[userInfo.Id]})
	}
	c.JSON(200, SearchUserResp{Users: users})
}
