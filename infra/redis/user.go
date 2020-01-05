package redis

import (
	"fmt"
	"strconv"
	"time"
)

const OnlineUser = "poseidon_online_user_"
const OnlineUserExpiration = 10 * time.Second

func AddUser(userId int64) error {
	return redisCli.Set(fmt.Sprintf("%s%d", OnlineUser, userId), "", OnlineUserExpiration).Err()
}

func GetUsers() ([]int64, error) {
	userIdsStr, err := redisCli.Keys(fmt.Sprintf("%s*", OnlineUser)).Result()
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, 0, len(userIdsStr))
	for _, userIdStr := range userIdsStr {
		userId, err := strconv.ParseInt(userIdStr[len(OnlineUser):], 10, 64)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func KickUser(userId int64) error {
	return redisCli.Del(fmt.Sprintf("%d", userId)).Err()
}
//
//func KickAllUser() error {
//	userIds, err := GetUsers()
//	if err != nil {
//		return err
//	}
//	userIdsStr := make([]string, 0, len(userIds))
//	for _, userId := range userIds {
//		userIdsStr = append(userIdsStr, fmt.Sprintf("%d", userId))
//	}
//	return redisCli.Del(userIdsStr...).Err()
//}
