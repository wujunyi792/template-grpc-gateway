package sonyflake

import (
	"github.com/sony/sonyflake"
	"pinnacle-primary-be/pkg/logger"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenSonyFlakeId() (int64, error) {
	id, err := flake.NextID()
	if err != nil {
		logger.NameSpace("sonyFlakeId").Warn("flake NextID failed: ", err)
		return 0, err
	}
	return int64(id), nil
}
