package entity

type User struct {
	Id             int64
	Password       string
	NickName       string
	CreateTime     int64
	LastOnlineTime int64
}
