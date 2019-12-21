package redis

import (
	"strconv"
)

const OnlineUserSet = "poseidon_online_user"

func AddUser(userId int64) error {
	return redisCli.SAdd(OnlineUserSet, userId).Err()
}

func GetUsers() ([]int64, error) {
	userIdsStr, err := redisCli.SMembers(OnlineUserSet).Result()
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, 0, len(userIdsStr))
	for _, userIdStr := range userIdsStr {
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func KickUser(userId int64) error {
	return redisCli.SRem(OnlineUserSet, userId).Err()
}

func KickAllUser() error {
	userIds, err := GetUsers()
	if err != nil {
		return err
	}
	userIdsInterface := make([]interface{}, 0, len(userIds))
	for _, userId := range userIds {
		userIdsInterface = append(userIdsInterface, userId)
	}
	return redisCli.SRem(OnlineUserSet, userIdsInterface...).Err()
}
