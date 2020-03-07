package gin_handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/mysql"
	"poseidon/infra/redis"
	"strconv"
	"time"
)

type Status struct {
	StatusCode    int
	StatusMessage string
}

type HeartBeatResp struct {
	Status
}

func HeartBeat(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(200, HeartBeatResp{Status: Status{StatusCode: 1, StatusMessage: err.Error()}})
		return
	}
	err = mysql.UpdateLastOnlineTime(userId, time.Now())
	if err != nil {
		c.JSON(200, HeartBeatResp{Status: Status{StatusCode: 1, StatusMessage: err.Error()}})
		return
	}
	err = redis.AddUser(userId)
	if err != nil {
		c.JSON(200, HeartBeatResp{Status: Status{StatusCode: 1, StatusMessage: err.Error()}})
		return
	}
	c.JSON(200, HeartBeatResp{})
}
