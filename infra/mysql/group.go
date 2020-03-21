package mysql

import (
	"poseidon/entity"
	"poseidon/utils"
	"time"
)

func CreateGroup(userId int64, name string) (int64, error) {
	id := utils.GenerateId(4)
	group := entity.Group{Id: id, Owner: userId, Name: name, CreateTime: time.Now().Unix()}
	return group.Id, db.Create(&group).Error
}

func GetGroupInfo(id int64) (*entity.Group, error) {
	var group entity.Group
	err := db.Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func GetGroupInfos(ids []int64) ([]*entity.Group, error) {
	var groups []*entity.Group
	err := db.Where("id IN (?)", ids).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func SearchGroup(data string) ([]*entity.Group, error) {
	var groups []*entity.Group
	res := db.Model(&entity.Group{}).Select("id, name, create_time").Where("id = ? OR name LIKE ?", data, "%"+data+"%").Find(&groups)
	if res.Error != nil {
		return nil, res.Error
	}
	return groups, nil
}

func DeleteGroup(groupId int64) error {
	err := DeleteGroupUser(groupId)
	if err != nil {
		return err
	}
	err = DeleteGroupMessage(groupId)
	if err != nil {
		return err
	}
	return db.Where("id = ?", groupId).Delete(entity.Group{}).Error
}
