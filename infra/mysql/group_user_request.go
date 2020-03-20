package mysql

import (
	"poseidon/entity"
	"time"
)

func CheckDuplicateGroupRequest(userId, groupId int64) (bool, error) {
	var count int
	if err := db.Model(&entity.GroupUserRequest{}).
		Where("(user_id_send = ? OR user_id_recv = ?) AND group_id = ? AND status = 0 AND parent_id = -1", userId, userId, groupId).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func CreateGroupUserRequest(userId, groupId int64, reqType entity.GroupUserRequestType) (*entity.GroupUserRequest, error) {
	groupUserRequest := entity.GroupUserRequest{UserIdSend: userId, GroupId: groupId, CreateTime: time.Now().Unix(), Status: entity.Pending, ParentId: -1, Type: reqType}
	return &groupUserRequest, db.Create(&groupUserRequest).Error
}

func AcceptedAddGroup(id int64) (userIdSend, groupId int64, err error) {
	tx := db.Begin()
	defer tx.Rollback()
	if err := tx.Model(&entity.GroupUserRequest{}).Where("id = ?", id).Update(map[string]interface{}{
		"status": entity.Accepted,
	}).Error; err != nil {
		return 0, 0, err
	}
	req := entity.GroupUserRequest{Id: id}
	if err := tx.Model(&req).Select("user_id_send, group_id").First(&req).Error; err != nil {
		return 0, 0, err
	}
	now := time.Now()
	if err := tx.Create(&entity.GroupUser{GroupId: req.GroupId, UserId: req.UserIdSend, CreateTime: now.Unix(), LastReadMsgId: -1}).Error; err != nil {
		return 0, 0, err
	}
	return req.UserIdSend, req.GroupId, tx.Commit().Error
}

func RejectedAddGroup(id int64) (userIdSend, groupId int64, err error) {
	if err := db.Model(&entity.GroupUserRequest{}).Where("id = ?", id).Update(map[string]interface{}{
		"status": entity.Rejected,
	}).Error; err != nil {
		return 0, 0, err
	}
	req := entity.GroupUserRequest{Id: id}
	if err := db.Model(&req).Select("user_id_send, group_id").First(&req).Error; err != nil {
		return 0, 0, err
	}
	return req.UserIdSend, req.GroupId, nil
}

func CreateReplyAddGroup(parentId, userIdSend, userIdRecv, groupId int64, now time.Time, reqType entity.GroupUserRequestType) (int64, error) {
	groupUserRequest := entity.GroupUserRequest{UserIdSend: userIdSend, UserIdRecv: userIdRecv, GroupId: groupId, CreateTime: now.Unix(), Status: entity.Pending, ParentId: parentId, Type: reqType}
	if err := db.Create(&groupUserRequest).Error; err != nil {
		return 0, err
	}
	return groupUserRequest.Id, nil
}

func SyncGroupUserRequest(userId, groupUserId int64) ([]*entity.GroupUserRequest, error) {
	groupIds, err := GetGroupList(userId)
	if err != nil {
		return nil, err
	}
	var groupUserRequests []*entity.GroupUserRequest
	ret := db.Table("group_user_request").Where("(user_id_recv = ? OR user_id_send = ? OR group_id IN (?)) AND id > ?", userId, userId, groupIds, groupUserId).Find(&groupUserRequests)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return groupUserRequests, nil
}

func GetGroupUserStatus(ids []int64) (map[int64]int32, error) {
	var res []*entity.GroupUserRequest
	ret := make(map[int64]int32)
	err := db.Model(&entity.GroupUserRequest{}).Select("id, status").Where("id in (?)", ids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	for _, obj := range res {
		ret[obj.Id] = int32(obj.Status)
	}
	return ret, nil
}
