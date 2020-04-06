package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

const OnlineUser = "poseidon_online_user_"
const AccessToken = "poseidon_access_token_"
const OnlineUserExpiration = 15 * time.Second

func AddUser(userId int64, accessToken string) error {
	var err error

	ret := redisCli.Get(fmt.Sprintf("%s%d", OnlineUser, userId))
	if ret.Err() != nil && ret.Err() != redis.Nil {
		return ret.Err()
	}
	if ret.Err() == nil {
		err = redisCli.Del(fmt.Sprintf("%s%s", AccessToken, ret.Val())).Err()
		if err != nil {
			return err
		}
	}

	err = redisCli.Set(fmt.Sprintf("%s%d", OnlineUser, userId), accessToken, OnlineUserExpiration).Err()
	if err != nil {
		return err
	}
	return redisCli.Set(fmt.Sprintf("%s%s", AccessToken, accessToken), userId, OnlineUserExpiration).Err()
}

func RefreshUser(userId int64, accessToken string) error {
	var err error
	err = redisCli.Set(fmt.Sprintf("%s%d", OnlineUser, userId), accessToken, OnlineUserExpiration).Err()
	if err != nil {
		return err
	}
	return redisCli.Set(fmt.Sprintf("%s%s", AccessToken, accessToken), userId, OnlineUserExpiration).Err()
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

func KickUser(userId int64, accessToken string) error {
	var err error
	err = redisCli.Del(fmt.Sprintf("%s%d", OnlineUser, userId)).Err()
	if err != nil {
		return err
	}
	return redisCli.Del(fmt.Sprintf("%s%s", AccessToken, accessToken)).Err()
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
