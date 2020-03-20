package entity

type GroupUser struct {
	Id            int64
	GroupId       int64
	UserId        int64
	CreateTime    int64
	LastReadMsgId int64
}
