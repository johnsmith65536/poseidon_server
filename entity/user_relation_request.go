package entity

type UserRelationRequest struct {
	Id         int64
	UserIdSend int64
	UserIdRecv int64
	CreateTime int64
	Status     Status
	ParentId   int64
}

type Status int32

const (
	Pending  Status = 0
	Rejected Status = 1
	Accepted Status = 2
)
