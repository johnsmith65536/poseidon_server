package entity

type UserRelation struct {
	Id            int64
	UserId        int64
	FriendUserId  int64
	CreateTime    int64
	LastReadMsgId int64
}
