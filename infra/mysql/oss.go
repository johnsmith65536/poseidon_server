package mysql

import "poseidon/entity"

func LoadSecretKey() (*entity.AccessKey, error) {
	var accessKey entity.AccessKey
	return &accessKey, db.Model(&entity.AccessKey{}).First(&accessKey).Error
}
