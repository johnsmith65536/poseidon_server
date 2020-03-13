package handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"poseidon/utils"
)

type LoginReq struct {
	UserId   int64
	Password string
}

type LoginResp struct {
	Success     bool
	AccessToken string
	Status
}

type LogoutReq struct {
	UserId      int64
	AccessToken string
}

type LogoutResp struct {
	Status
}

func Login(c *gin.Context) {
	var err error
	var req LoginReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	ok, err := mysql.Login(req.UserId, req.Password)
	PanicIfError(err)
	if !ok {
		c.JSON(200, LoginResp{Success: false})
		return
	}
	accessToken := utils.GenerateToken(10)
	err = redis.AddUser(req.UserId, accessToken)
	PanicIfError(err)
	c.JSON(200, LoginResp{Success: true, AccessToken: accessToken})
}

func Logout(c *gin.Context) {
	var err error
	var req LogoutReq
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	err = redis.KickUser(req.UserId, req.AccessToken)
	PanicIfError(err)
	c.JSON(200, LogoutResp{})
}
