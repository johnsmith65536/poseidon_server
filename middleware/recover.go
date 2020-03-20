package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"runtime"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				logrus.Errorf("%+v", err)
				PrintStack()
				c.JSON(200, map[string]interface{}{"StatusCode": 255, "StatusMessage": err.Error()})
			}
		}()
		c.Next()
	}
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}
