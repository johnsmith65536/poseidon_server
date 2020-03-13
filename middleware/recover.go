package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				logrus.Error(err)
				c.JSON(200, map[string]interface{}{"StatusCode": 255, "StatusMessage": err.Error()})
			}
		}()
		c.Next()
	}
}
