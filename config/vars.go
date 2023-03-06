package config

import (
	"pinnacle-primary-be/core/discov"
	"pinnacle-primary-be/core/sentryx"
	"pinnacle-primary-be/core/store/mysql"
	"pinnacle-primary-be/core/store/rds"
	"pinnacle-primary-be/core/tracex"
)

type GlobalConfig struct {
	MODE        string          `yaml:"Mode"`
	ProgramName string          `yaml:"ProgramName"`
	BaseURL     string          `yaml:"BaseURL"`
	AUTHOR      string          `yaml:"Author"`
	Listen      string          `yaml:"Listen"`
	Port        string          `yaml:"Port"`
	MainMysql   mysql.OrmConf   `yaml:"MainMysql"`
	MainCache   rds.RedisConf   `yaml:"MainCache"`
	ETCD        discov.EtcdConf `yaml:"Etcd"`
	Trace       tracex.Config   `yaml:"Trace"`
	Sentry      sentryx.Config  `yaml:"Sentry"`
	Auth        struct {
		Secret string `yaml:"Secret"`
		Issuer string `yaml:"Issuer"`
	} `yaml:"Auth"`
}
