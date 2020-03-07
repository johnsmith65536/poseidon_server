package handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/mysql"
)

type CreateObjectReq struct {
	ETag string
	Name string
}

type CreateObjectResp struct {
	Id int64
	Status
}

func CreateObject(c *gin.Context) {
	var err error
	var req CreateObjectReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, CreateObjectResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}

	id, err := mysql.CreateObject(req.ETag, req.Name)
	if err != nil {
		c.JSON(200, CreateObjectResp{Status: Status{StatusCode: 255, StatusMessage: err.Error()}})
		return
	}
	c.JSON(200, CreateObjectResp{Id: id})
}
