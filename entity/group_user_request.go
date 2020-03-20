package entity

type GroupUserRequest struct {
	Id         int64
	UserIdSend int64
	UserIdRecv int64
	GroupId    int64
	CreateTime int64
	Status     Status
	ParentId   int64
	Type       GroupUserRequestType
}

type GroupUserRequestType int32

const (
	AddGroup       GroupUserRequestType = 0
	InviteAddGroup GroupUserRequestType = 1
)
