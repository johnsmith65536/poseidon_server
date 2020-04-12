package handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/entity"
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
	UserId int64
	Data string
}*/

type SearchUserResp struct {
	Users []*User
	Status
}

/*type GetUserInfoReq struct {
	UserId int64
}*/

type GetUserInfoResp struct {
	User *entity.User
	Status
}

func CreateUser(c *gin.Context) {
	var err error
	var req CreateUserReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)

	userId := utils.GenerateId(4)
	err = mysql.CreateUser(userId, req.Password, req.NickName, time.Now())
	PanicIfError(err)
	c.JSON(200, CreateUserResp{UserId: userId})
}

func SearchUser(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	PanicIfError(err)

	data := c.Query("data")

	users := make([]*User, 0)
	userIds := make([]int64, 0)
	userInfos, err := mysql.SearchUser(data)
	PanicIfError(err)
	for _, userInfo := range userInfos {
		userIds = append(userIds, userInfo.Id)
	}
	friendUserIds, err := mysql.GetRelation(userId, userIds)
	PanicIfError(err)
	friendMap := make(map[int64]bool)
	for _, friendUserId := range friendUserIds {
		friendMap[friendUserId] = true
	}
	for _, userInfo := range userInfos {
		users = append(users, &User{Id: userInfo.Id, NickName: userInfo.NickName, LastOnlineTime: userInfo.LastOnlineTime, IsFriend: friendMap[userInfo.Id]})
	}
	c.JSON(200, SearchUserResp{Users: users})
}

func GetUserInfo(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	PanicIfError(err)

	res, err := mysql.GetUserNickNames([]int64{userId})
	PanicIfError(err)

	c.JSON(200, GetUserInfoResp{User: &entity.User{Id: userId, NickName: res[userId]}})
}
