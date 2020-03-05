package mysql

import (
	"poseidon/entity"
	"poseidon/thrift"
)

func CreateObject(eTag string, name string) (int64, error) {
	var object entity.Object
	err := db.FirstOrCreate(&object, entity.Object{ETag: eTag, Name: name}).Error
	if err != nil {
		return 0, err
	}
	return object.Id, nil
}

func SyncObject(objIds []int64) ([]*thrift.Object, error) {
	var objects []*thrift.Object
	err := db.Where("id IN (?)", objIds).Find(&objects).Error
	if err != nil {
		return nil, err
	}
	return objects, nil
}

func GetObject(objId int64) (*entity.Object, error) {
	var object entity.Object
	err := db.Where("id = ?", objId).First(&object).Error
	if err != nil {
		return nil, err
	}
	return &object, nil
}
