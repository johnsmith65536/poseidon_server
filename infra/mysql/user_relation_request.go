package mysql

import (
	"poseidon/entity"
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
	maxMsgId, err := GetFriendLastMsgId(req.UserIdSend, req.UserIdRecv)
	if err != nil {
		return 0, 0, err
	}
	now := time.Now()
	if err := tx.Create(&entity.UserRelation{UserId: req.UserIdSend, FriendUserId: req.UserIdRecv, CreateTime: now.Unix(), LastReadMsgId: maxMsgId}).Error; err != nil {
		return 0, 0, err
	}
	if err := tx.Create(&entity.UserRelation{UserId: req.UserIdRecv, FriendUserId: req.UserIdSend, CreateTime: now.Unix(), LastReadMsgId: maxMsgId}).Error; err != nil {
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

func SyncUserRelationRequest(userId, userRelationId int64) ([]*entity.UserRelationRequest, error) {
	var userRelationRequests []*entity.UserRelationRequest
	ret := db.Table("user_relation_request").Where("(user_id_recv = ? OR user_id_send = ?) AND id > ?", userId, userId, userRelationId).Find(&userRelationRequests)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return userRelationRequests, nil
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

func UpdateUserRelationRequestStatus(userRelationRequestIds map[int64]int32) error {
	for id, val := range userRelationRequestIds {
		err := db.Model(&entity.UserRelationRequest{}).Where("id = ?", id).Update(map[string]interface{}{
			"status": val,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
