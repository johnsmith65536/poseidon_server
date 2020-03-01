package redis

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

const messageChannel = "poseidon_message_channel_"

type BroadcastMsgType int

const (
	Chat           BroadcastMsgType = 0
	AddFriend      BroadcastMsgType = 1
	ReplyAddFriend BroadcastMsgType = 2
)

func BroadcastMessage(userId int64, data map[string]interface{}, msgType BroadcastMsgType) error {
	data["BroadcastMsgType"] = msgType
	msg, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Marshal failed, err: %+v", err)
		return err
	}
	if err = redisCli.Publish(fmt.Sprintf("%s%d", messageChannel, userId), string(msg)).Err(); err != nil {
		logrus.Errorf("Publish failed, err: %+v", err)
	}
	return err
}
