package mysql

import (
	"poseidon/entity"
	"time"
)

func WriteMessage(userIdSend, userIdRecv int64, groupId int64, byteContent []byte, createTime time.Time, contentType int32, msgType int32) (*entity.Message, error) {
	message := entity.Message{UserIdSend: userIdSend, UserIdRecv: userIdRecv, GroupId: groupId, Content: byteContent, CreateTime: createTime.Unix(), ContentType: contentType, MsgType: msgType}
	return &message, db.Create(&message).Error
}

func SyncMessage(userId, messageId int64) ([]*entity.Message, error) {
	var messages []*entity.Message
	groupIds, err := GetGroupList(userId)
	if err != nil {
		return nil, err
	}
	ret := db.Model(&entity.Message{}).Where("(user_id_send = ? OR user_id_recv = ? OR group_id IN (?)) AND id > ?", userId, userId, groupIds, messageId).Find(&messages)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return messages, nil
}

func DeleteGroupMessage(groupId int64) error {
	return db.Where("group_id = ?", groupId).Delete(entity.Message{}).Error
}

func FetchFriendHistoryMessage(userIdAlice, userIdBob, localCount int64) ([]*entity.Message, error) {
	var count int64
	messages := make([]*entity.Message, 0)
	sql := "(user_id_send = ? AND user_id_recv = ?) OR (user_id_send = ? AND user_id_recv = ?)"
	ret := db.Model(&entity.Message{}).Where(sql, userIdAlice, userIdBob, userIdBob, userIdAlice).Count(&count)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if localCount == count {
		return messages, nil
	}
	ret = db.Model(&entity.Message{}).Where(sql, userIdAlice, userIdBob, userIdBob, userIdAlice).Find(&messages)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return messages, nil
}

func FetchGroupHistoryMessage(groupId, localCount int64) ([]*entity.Message, error) {
	var count int64
	messages := make([]*entity.Message, 0)
	sql := "group_id = ?"
	ret := db.Model(&entity.Message{}).Where(sql, groupId).Count(&count)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if localCount == count {
		return messages, nil
	}
	ret = db.Model(&entity.Message{}).Where(sql, groupId).Find(&messages)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return messages, nil
}

func GetGroupLastMsgId(groupId int64) (int64, error) {
	var message entity.Message
	ret := db.Model(&entity.Message{}).Select("max(id) as id").Where("group_id = ?", groupId).First(&message)
	if ret.Error != nil && !ret.RecordNotFound() {
		return 0, ret.Error
	}
	if ret.RecordNotFound() {
		return -1, nil
	}
	return message.Id, nil
}

func GetFriendLastMsgId(userId, friendUserId int64) (int64, error) {
	var message entity.Message
	ret := db.Model(&entity.Message{}).Select("max(id) as id").Where("(user_id_send = ? AND user_id_recv = ?) OR (user_id_send = ? AND user_id_recv = ?)", userId, friendUserId, friendUserId, userId).First(&message)
	if ret.Error != nil && !ret.RecordNotFound() {
		return 0, ret.Error
	}
	if ret.RecordNotFound() {
		return -1, nil
	}
	return message.Id, nil
}
