package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"strconv"
	"time"
)


type HeartBeatResp struct {
	Status
}

func HeartBeat(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	PanicIfError(err)
	err = mysql.UpdateLastOnlineTime(userId, time.Now())
	PanicIfError(err)
	err = redis.RefreshUser(userId, c.GetHeader("access_token"))
	PanicIfError(err)
	c.JSON(200, HeartBeatResp{})
}

func Ping(c *gin.Context) {
	fmt.Println(c.QueryArray("data"))
	fmt.Println(c.QueryArray("info"))
}
