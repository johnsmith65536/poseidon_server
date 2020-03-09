package redis

import (
	"fmt"
)

func CheckAccessToken(accessToken string) (bool, error) {
	accessTokens, err := redisCli.Keys(fmt.Sprintf("%s%s", AccessToken, accessToken)).Result()
	if err != nil {
		return false, err
	}
	return len(accessTokens) == 1, nil
}
