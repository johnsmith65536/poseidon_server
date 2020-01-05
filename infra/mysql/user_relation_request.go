package mysql

import (
	"poseidon/entity"
	"poseidon/thrift"
	"time"
)

func CreateUserRelationRequest(userIdSend, userIdRecv int64) (*entity.UserRelationRequest, error) {
	userRelationRequest := entity.UserRelationRequest{UserIdSend: userIdSend, UserIdRecv: userIdRecv, CreateTime: time.Now().Unix(), Status: entity.Pending}
	return &userRelationRequest, db.Create(&userRelationRequest).Error
}

func AcceptedAddFriend(id int64) error {
	tx := db.Begin()
	defer tx.Rollback()
	if err := tx.Model(&entity.UserRelationRequest{}).Where("id = ?", id).Update(map[string]interface{}{
		"status": entity.Accepted,
	}).Error; err != nil {
		return err
	}
	req := entity.UserRelationRequest{Id: id}
	if err := tx.Model(&req).Select("user_id_send, user_id_recv").First(&req).Error; err != nil {
		return err
	}
	now := time.Now()
	if err := tx.Create(&entity.UserRelation{UserId: req.UserIdSend, FriendUserId: req.UserIdRecv, CreateTime: now.Unix()}).Error; err != nil {
		return err
	}
	if err := tx.Create(&entity.UserRelation{UserId: req.UserIdRecv, FriendUserId: req.UserIdSend, CreateTime: now.Unix()}).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

func RejectedAddFriend(id int64) error {
	return db.Model(&entity.UserRelationRequest{}).Where("id = ?", id).Update(map[string]interface{}{
		"status": entity.Rejected,
	}).Error
}

func FetchOfflineUserRelation(userId, userRelationId int64) ([]*thrift.UserRelation, error) {
	var userRelations []*thrift.UserRelation
	ret := db.Model(&entity.UserRelationRequest{}).Where("user_id_recv = ? AND id = ?", userId, userRelationId).Find(&userRelations)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return userRelations, nil
}
