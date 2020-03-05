package mysql

import (
	"poseidon/entity"
	"poseidon/thrift"
	"time"
)

func CreateUserRelationRequest(userIdSend, userIdRecv int64) (*entity.UserRelationRequest, error) {
	userRelationRequest := entity.UserRelationRequest{UserIdSend: userIdSend, UserIdRecv: userIdRecv, CreateTime: time.Now().Unix(), Status: entity.Pending, ParentId: -1}
	return &userRelationRequest, db.Create(&userRelationRequest).Error
}

func AcceptedAddFriend(id int64) (userIdSend int64, userIdRecv int64, err error) {
	tx := db.Begin()
	defer tx.Rollback()
	if err := tx.Model(&entity.UserRelationRequest{}).Where("id = ?", id).Update(map[string]interface{}{
		"status": entity.Accepted,
	}).Error; err != nil {
		return 0, 0, err
	}
	req := entity.UserRelationRequest{Id: id}
	if err := tx.Model(&req).Select("user_id_send, user_id_recv").First(&req).Error; err != nil {
		return 0, 0, err
	}
	now := time.Now()
	if err := tx.Create(&entity.UserRelation{UserId: req.UserIdSend, FriendUserId: req.UserIdRecv, CreateTime: now.Unix()}).Error; err != nil {
		return 0, 0, err
	}
	if err := tx.Create(&entity.UserRelation{UserId: req.UserIdRecv, FriendUserId: req.UserIdSend, CreateTime: now.Unix()}).Error; err != nil {
		return 0, 0, err
	}
	return req.UserIdSend, req.UserIdRecv, tx.Commit().Error
}

func RejectedAddFriend(id int64) (userIdSend int64, userIdRecv int64, err error) {
	if err := db.Model(&entity.UserRelationRequest{}).Where("id = ?", id).Update(map[string]interface{}{
		"status": entity.Rejected,
	}).Error; err != nil {
		return 0, 0, err
	}
	req := entity.UserRelationRequest{Id: id}
	if err := db.Model(&req).Select("user_id_send, user_id_recv").First(&req).Error; err != nil {
		return 0, 0, err
	}
	return req.UserIdSend, req.UserIdRecv, nil
}

func SyncUserRelation(userId, userRelationId int64) ([]*thrift.UserRelation, error) {
	var userRelations []*thrift.UserRelation
	ret := db.Table("user_relation_request").Where("(user_id_recv = ? OR user_id_send = ?) AND id > ?", userId, userId, userRelationId).Find(&userRelations)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return userRelations, nil
}

func CreateReplyAddFriend(parentId int64, userIdSend int64, userIdRecv int64, now time.Time) (int64, error) {
	userRelationRequest := entity.UserRelationRequest{UserIdSend: userIdSend, UserIdRecv: userIdRecv, CreateTime: now.Unix(), Status: entity.Pending, ParentId: parentId}
	if err := db.Create(&userRelationRequest).Error; err != nil {
		return 0, err
	}
	return userRelationRequest.Id, nil
}

func CheckDuplicateRequest(userIdSend int64, userIdRecv int64) (bool, error) {
	var count int
	if err := db.Model(&entity.UserRelationRequest{}).
		Where("((user_id_send = ? AND user_id_recv = ?) OR (user_id_send = ? AND user_id_recv = ?)) AND status = 0 AND parent_id = -1", userIdSend, userIdRecv, userIdRecv, userIdSend).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func GetRelationStatus(ids []int64) (map[int64]int32, error) {
	var res []*entity.UserRelationRequest
	ret := make(map[int64]int32)
	err := db.Model(&entity.UserRelationRequest{}).Select("id, status").Where("id in (?)", ids).Find(&res).Error
	if err != nil {
		return nil, err
	}
	for _, obj := range res {
		ret[obj.Id] = int32(obj.Status)
	}
	return ret, nil
}
