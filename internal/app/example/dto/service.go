package dto

import (
	"errors"
	"gorm.io/gorm"
	"pinnacle-primary-be/internal/app/example/model"
	"pinnacle-primary-be/internal/database"
)

func init() {
	// 就近原则 那里要用 哪里就马上初始化
	_ = database.AutoMigrate("*", &model.Example{})
}

var Example *services

type services struct{}

func (*services) Count() (count int64) {
	database.GetDB("*").Model(&model.Example{}).Count(&count)
	return
}

func (*services) Save(data *model.Example, tx ...*gorm.DB) error {
	if len(tx) != 0 {
		return tx[0].Save(data).Error
	}
	return database.GetDB("*").Save(data).Error
}

func (*services) Find(where *model.Example) ([]*model.Example, error) {
	list := make([]*model.Example, 0, 10)
	err := database.GetDB("*").Where(where).Find(&list).Error
	return list, err
}

func (*services) Update(data *model.Example, tx ...*gorm.DB) error {
	var result *gorm.DB
	if len(tx) != 0 {
		result = tx[0].Model(data).Updates(data)
	} else {
		result = database.GetDB("*").Model(data).Updates(data)
	}
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected != 1 {
		return errors.New("update locked by optimistic lock")
	}
	return nil
}

func (*services) Del(data *model.Example, tx ...*gorm.DB) error {
	if len(tx) != 0 {
		return tx[0].Where(data).Delete(&model.Example{}).Error
	}
	return database.GetDB("*").Where(data).Delete(&model.Example{}).Error
}

func (*services) GetByID(sid uint) (s *model.Example, err error) {
	s = &model.Example{}
	err = database.GetDB("*").Where("id = ?", sid).First(s).Error
	return
}

func (*services) Gets(offset, limit int) ([]*model.Example, error) {
	var users []*model.Example
	err := database.GetDB("*").Model(&model.Example{}).Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}
