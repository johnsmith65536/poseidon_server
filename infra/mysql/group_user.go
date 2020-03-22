package mysql

import (
	"poseidon/entity"
	"time"
)

func AddGroupMember(groupId int64, userIds []int64) error {
	for _, userId := range userIds {
		if err := db.Create(&entity.GroupUser{GroupId: groupId, UserId: userId, CreateTime: time.Now().Unix(), LastReadMsgId: -1}).Error; err != nil {
			return err
		}
	}
	return nil
}

func CheckIsGroupMember(userId int64, groupId int64) (bool, error) {
	var count int
	err := db.Model(&entity.GroupUser{}).Where("group_id = ? AND user_id = ?", groupId, userId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func GetGroupList(userId int64) ([]int64, error) {
	groupIds := make([]int64, 0)
	return groupIds, db.Model(&entity.GroupUser{}).Select("group_id").Where("user_id = ?", userId).Pluck("group_id", &groupIds).Error
}

func GetMemberList(groupId int64) ([]int64, error) {
	var userIds []int64
	return userIds, db.Model(&entity.GroupUser{}).Select("user_id").Where("group_id = ?", groupId).Pluck("user_id", &userIds).Error
}

func GetGroupLastReadMsgId(userId int64) (map[int64]int64, error) {
	ret := make(map[int64]int64)
	var groupUsers []*entity.GroupUser
	if err := db.Model(&entity.GroupUser{}).Select("group_id, last_read_msg_id").Where("user_id = ?", userId).Find(&groupUsers).Error; err != nil {
		return nil, err
	}
	for _, groupUser := range groupUsers {
		ret[groupUser.GroupId] = groupUser.LastReadMsgId
	}
	return ret, nil
}

func UpdateGroupLastReadMsgId(userId int64, lastReadMsgId map[int64]int64) error {
	for groupId, msgId := range lastReadMsgId {
		if err := db.Model(&entity.GroupUser{}).Where("group_id = ? AND user_id = ?", groupId, userId).Update(map[string]interface{}{
			"last_read_msg_id": msgId,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func DeleteGroupUser(groupId int64) error {
	return db.Where("group_id = ?", groupId).Delete(entity.GroupUser{}).Error
}

func DeleteGroupMember(groupId, userId int64) error {
	return db.Where("group_id = ? AND user_id = ?", groupId, userId).Delete(entity.GroupUser{}).Error
}
