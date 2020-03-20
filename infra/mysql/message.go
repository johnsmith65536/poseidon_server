package mysql

import (
	"poseidon/entity"
	"time"
)

func WriteMessage(userIdSend, userIdRecv int64, groupId int64, byteContent []byte, createTime time.Time, contentType int32, msgType int32, isRead bool) (*entity.Message, error) {
	message := entity.Message{UserIdSend: userIdSend, UserIdRecv: userIdRecv, GroupId: groupId, Content: byteContent, CreateTime: createTime.Unix(), ContentType: contentType, MsgType: msgType, IsRead: isRead}
	return &message, db.Create(&message).Error
}

func UpdateMessageStatus(messageIds, userRelationRequestIds, groupUserRequestIds map[int64]int32, ) error {
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

	for id, val := range groupUserRequestIds {
		err := db.Model(&entity.GroupUserRequest{}).Where("id = ?", id).Update(map[string]interface{}{
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
