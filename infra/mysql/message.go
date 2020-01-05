package mysql

import (
	"poseidon/entity"
	"poseidon/thrift"
	"time"
)

func WriteMessage(userIdSend, userIdRecv int64, groupId int64, content string, createTime time.Time, msgType int32, isRead bool) (*entity.Message, error) {
	message := entity.Message{UserIdSend: userIdSend, UserIdRecv: userIdRecv, GroupId: groupId, Content: content, CreateTime: createTime.Unix(), MsgType: msgType, IsRead: isRead}
	return &message, db.Create(&message).Error
}

func MessageDelivered(msgId int64) error {
	return db.Model(&entity.Message{}).Where("id = ?", msgId).Update(map[string]interface{}{
		"delivered": true,
	}).Error
}

func FetchOfflineMessage(userId, messageId int64) ([]*thrift.Message, error) {
	var messages []*thrift.Message
	ret := db.Model(&entity.Message{}).Where("user_id_recv = ? AND id > ?", userId, messageId).Find(&messages)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return messages, nil
}
