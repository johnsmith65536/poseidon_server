package middleware

import (
	"github.com/gin-gonic/gin"
	"poseidon/infra/redis"
)

func Auth(c *gin.Context) {
	accessToken := c.GetHeader("access_token")
	ok, err := redis.CheckAccessToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(200, map[string]interface{}{"StatusCode": 255, "StatusMessage": err.Error()})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(200, map[string]interface{}{"StatusCode": 254, "StatusMessage": "login status expired"})
		return
	}
}
