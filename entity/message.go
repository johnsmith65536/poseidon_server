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
}

type ContentType int32

const (
	Text       ContentType = 0
	ObjectData ContentType = 1
	Vibration  ContentType = 2
	Image      ContentType = 3
)

type MsgType int32

const (
	PrivateChat MsgType = 0
	GroupChat   MsgType = 1
)
