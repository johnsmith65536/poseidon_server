package mysql

import (
	"poseidon/entity"
	"time"
)

func CreateUser(userId int64, password, nickName string, time time.Time) error {
	user := entity.User{Id: userId, Password: password, NickName: nickName, CreateTime: time, LastOnlineTime: time}
	return db.Create(&user).Error
}

func Login(userId int64, password string) (bool, error) {
	var count int
	res := db.Model(&entity.User{}).Where("id = ? AND password = ?", userId, password).Count(&count)
	if res.Error != nil {
		return false, res.Error
	} else if count == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
