package mysql

import (
	"poseidon/entity"
)

func GetFriendsList(userId int64) ([]int64, error) {
	var userIds []int64
	err := db.Model(&entity.UserRelation{}).Select("friend_user_id").Where("user_id = ?", userId).Pluck("friend_user_id", &userIds).Error
	return userIds, err
}

func GetRelation(userId int64, userIds []int64) ([]int64, error) {
	var friendUserIds []int64
	err := db.Model(&entity.UserRelation{}).Select("friend_user_id").Where("user_id = ? AND friend_user_id IN (?)", userId, userIds).Pluck("friend_user_id", &friendUserIds).Error
	return friendUserIds, err
}

func DeleteFriend(userIdSend int64, userIdRecv int64) error {
	return db.Where("(user_id = ? AND friend_user_id = ?) OR (user_id = ? AND friend_user_id = ?)", userIdSend, userIdRecv, userIdRecv, userIdSend).Delete(entity.UserRelation{}).Error
}

func CheckIsFriend(userIdSend int64, userIdRecv int64) (bool, error) {
	var count int
	err := db.Model(&entity.UserRelation{}).Where("user_id = ? AND friend_user_id = ?", userIdSend, userIdRecv).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func GetFriendLastReadMsgId(userId int64) (map[int64]int64, error) {
	ret := make(map[int64]int64)
	var userRelations []*entity.UserRelation
	if err := db.Model(&entity.UserRelation{}).Select("friend_user_id, last_read_msg_id").Where("user_id = ?", userId).Find(&userRelations).Error; err != nil {
		return nil, err
	}
	for _, userRelation := range userRelations {
		ret[userRelation.FriendUserId] = userRelation.LastReadMsgId
	}
	return ret, nil
}

func UpdateFriendLastReadMsgId(userId int64, lastReadMsgId map[int64]int64) error {
	for friendUserId, msgId := range lastReadMsgId {
		if err := db.Model(&entity.UserRelation{}).Where("user_id = ? AND friend_user_id = ?", userId, friendUserId).Update(map[string]interface{}{
			"last_read_msg_id": msgId,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}