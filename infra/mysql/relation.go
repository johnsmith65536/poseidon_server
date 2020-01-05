package mysql

import (
	"poseidon/entity"
)

func GetFriendsList(userId int64) ([]int64, error) {
	var userIds []int64
	err := db.Model(&entity.UserRelation{}).Select("friend_user_id").Where("user_id = ?", userId).Pluck("friend_user_id", &userIds).Error
	return userIds, err
}

