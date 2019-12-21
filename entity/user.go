package entity

import "time"

type User struct {
	Id             int64
	Password       string
	NickName       string
	CreateTime     time.Time
	LastOnlineTime time.Time
}
