package handler

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/mysql"
	"strconv"
)

type UpdateRequestStatusReq struct {
	UserRelationRequestIds map[int64]int32
	GroupUserRequestIds    map[int64]int32
}

type UpdateRequestStatusResp struct {
	Status
}

/*type FetchRequestStatusReq struct {
	UserRelationRequestIds []int64
	GroupUserRequestIds []int64
}*/

type FetchRequestStatusResp struct {
	UserRelationRequestIds map[int64]int32
	GroupUserRequestIds    map[int64]int32
	Status
}

func UpdateRequestStatus(c *gin.Context) {
	var req UpdateRequestStatusReq
	var err error
	err = c.ShouldBindJSON(&req)
	PanicIfError(err)
	err = mysql.UpdateUserRelationRequestStatus(req.UserRelationRequestIds)
	PanicIfError(err)
	err = mysql.UpdateGroupUserRequestStatus(req.GroupUserRequestIds)
	PanicIfError(err)
	c.JSON(200, UpdateRequestStatusResp{})
}

func FetchRequestStatus(c *gin.Context) {
	var err error
	var userRelationRequestIds []int64
	var groupUserRequestIds []int64

	userRelationRequestIdsStr := c.QueryArray("user_relation_request_ids")
	for _, userRelationRequestIdStr := range userRelationRequestIdsStr {
		userRelationRequestId, err := strconv.ParseInt(userRelationRequestIdStr, 10, 64)
		PanicIfError(err)
		userRelationRequestIds = append(userRelationRequestIds, userRelationRequestId)
	}

	groupUserRequestIdsStr := c.QueryArray("group_user_request_ids")
	for _, groupUserRequestIdStr := range groupUserRequestIdsStr {
		groupUserRequestId, err := strconv.ParseInt(groupUserRequestIdStr, 10, 64)
		PanicIfError(err)
		groupUserRequestIds = append(groupUserRequestIds, groupUserRequestId)
	}

	relationStatus, err := mysql.GetRelationStatus(userRelationRequestIds)
	PanicIfError(err)
	groupUserStatus, err := mysql.GetGroupUserStatus(groupUserRequestIds)
	PanicIfError(err)
	c.JSON(200, FetchRequestStatusResp{UserRelationRequestIds: relationStatus, GroupUserRequestIds: groupUserStatus})
}
