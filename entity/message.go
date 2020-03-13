package entity

type Message struct {
	Id          int64
	UserIdSend  int64
	UserIdRecv  int64
	GroupId     int64
	Content     []byte
	CreateTime  int64
	ContentType int32
	MsgType     int32
	IsRead      bool
}

type ContentType int32

const (
	Text       ContentType = 0
	ObjectData ContentType = 1
	Vibration  ContentType = 2
)
