package mysql

import (
	"errors"
	"poseidon/entity"
	"time"
)

func CreateUser(userId int64, password, nickName string, time time.Time) error {
	user := entity.User{Id: userId, Password: password, NickName: nickName, CreateTime: time.Unix(), LastOnlineTime: time.Unix()}
	return db.Create(&user).Error
}

func Login(userId int64, password string) (bool, string, error) {
	var user entity.User
	res := db.Model(&entity.User{}).Select("nick_name").Where("id = ? AND password = ?", userId, password).First(&user)
	if res.Error != nil {
		if res.RecordNotFound() {
			return false, "", nil
		}
		return false, "", res.Error
	}
	return true, user.NickName, nil
}

func UpdateLastOnlineTime(userId int64, time time.Time) error {
	return db.Model(&entity.User{}).Where("id = ?", userId).Update(map[string]interface{}{
		"last_online_time": time.Unix(),
	}).Error
}

func GetUserNickNames(userIds []int64) (map[int64]string, error) {
	ret := make(map[int64]string)
	var users []*entity.User
	res := db.Model(&entity.User{}).Select("id, nick_name").Where("id IN (?)", userIds).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(userIds) != len(users) {
		return nil, errors.New("length not equal")
	}
	for _, user := range users {
		ret[user.Id] = user.NickName
	}
	return ret, nil
}

func SearchUser(data string) ([]entity.User, error) {
	var users []entity.User
	res := db.Model(&entity.User{}).Select("id, nick_name, last_online_time").Where("id = ? OR nick_name LIKE ?", data, "%"+data+"%").Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}
