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

func UpdateMessageStatus(messageIds map[int64]int32, userRelationRequestIds map[int64]int32) error {
	for id, val := range messageIds {
		var isRead bool
		if val == 1 {
			isRead = true
		}
		err := db.Model(&entity.Message{}).Where("id = ?", id).Update(map[string]interface{}{
			"is_read": isRead,
		}).Error
		if err != nil {
			return err
		}
	}

	for id, val := range userRelationRequestIds {
		err := db.Model(&entity.UserRelationRequest{}).Where("id = ?", id).Update(map[string]interface{}{
			"status": val,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func SyncMessage(userId, messageId int64) ([]*thrift.Message, error) {
	var messages []*thrift.Message
	ret := db.Model(&entity.Message{}).Where("user_id_recv = ? AND id > ?", userId, messageId).Find(&messages)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return messages, nil
}
