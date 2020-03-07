package handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/oss"
	"strconv"
)

/*type GetSTSInfoReq struct {
	UserId int64 `thrift:"userId,1" db:"userId" json:"userId"`
}*/

type GetSTSInfoResp struct {
	SecurityToken   string
	AccessKeyId     string
	AccessKeySecret string
	Status
}

func GetSTSInfo(c *gin.Context) {
	var err error
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(200, GetSTSInfoResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}

	info, err := oss.GetSTSInfo(userId)
	if err != nil {
		if err != nil {
			c.JSON(200, GetSTSInfoResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
			return
		}
	}
	c.JSON(200, GetSTSInfoResp{SecurityToken: info.SecurityToken, AccessKeyId: info.AccessKeyId, AccessKeySecret: info.AccessKeySecret})
}
