package mysql

import (
	"poseidon/entity"
	"time"
)

func WriteMessage(userIdSend, userIdRecv int64, groupId int64, content string, createTime time.Time, contentType int32, msgType int32, isRead bool) (*entity.Message, error) {
	message := entity.Message{UserIdSend: userIdSend, UserIdRecv: userIdRecv, GroupId: groupId, Content: content, CreateTime: createTime.Unix(), ContentType: contentType, MsgType: msgType, IsRead: isRead}
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

func SyncMessage(userId, messageId int64) ([]*entity.Message, error) {
	var messages []*entity.Message
	ret := db.Model(&entity.Message{}).Where("(user_id_send = ? OR user_id_recv = ?) AND id > ?", userId, userId, messageId).Find(&messages)
	if ret.Error != nil && !ret.RecordNotFound() {
		return nil, ret.Error
	}
	return messages, nil
}

func GetMessageStatus(ids []int64) (map[int64]int32, error) {
	var messages []entity.Message
	res := make(map[int64]int32)
	err := db.Model(&entity.Message{}).Select("id, is_read").Where("id in (?)", ids).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	for _, message := range messages {
		if message.IsRead {
			res[message.Id] = 1
		} else {
			res[message.Id] = 0
		}
	}
	return res, nil
}
