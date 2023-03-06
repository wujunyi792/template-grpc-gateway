package dao

import (
	"pinnacle-primary-be/core/store/mysql"
	"pinnacle-primary-be/internal/app/ping/model"
)

var (
	Ping *mysql.Orm
)

func AutoMigrate() error {
	return Ping.AutoMigrate(&model.Ping{})
}
