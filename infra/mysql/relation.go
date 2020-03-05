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
	return count == 0, nil
}
