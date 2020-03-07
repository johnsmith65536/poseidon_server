package handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
)

type LoginReq struct {
	UserId   int64
	Password string
}

type LoginResp struct {
	Success bool
	Status
}

type LogoutReq struct {
	UserId int64
}

type LogoutResp struct {
	Status
}

func Login(c *gin.Context) {
	var err error
	var req LoginReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, LoginResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	ok, err := mysql.Login(req.UserId, req.Password)
	if err != nil {
		c.JSON(200, LoginResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	if !ok {
		c.JSON(200, LoginResp{Success: false})
		return
	}
	err = redis.AddUser(req.UserId)
	if err != nil {
		c.JSON(200, LoginResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	c.JSON(200, LoginResp{Success: true})
}

func Logout(c *gin.Context) {
	var err error
	var req LogoutReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, LogoutResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	err = redis.KickUser(req.UserId)
	if err != nil {
		c.JSON(200, LogoutResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	c.JSON(200, LogoutResp{})
}
