package entity

type Message struct {
	Id         int64
	UserIdSend int64
	UserIdRecv int64
	GroupId    int64
	Content    string
	CreateTime int64
	MsgType    int32
	IsRead     bool
}
