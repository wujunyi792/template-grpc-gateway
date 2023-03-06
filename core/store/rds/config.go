package rds

import "errors"

var (
	// ErrEmptyHost is an error that indicates no redis host is set.
	ErrEmptyHost = errors.New("empty redis host")
	// ErrEmptyType is an error that indicates no redis type is set.
	ErrEmptyType = errors.New("empty redis type")
	// ErrEmptyKey is an error that indicates no redis key is set.
	ErrEmptyKey = errors.New("empty redis key")
	// ErrPing is an error that indicates ping failed.
	ErrPing = errors.New("ping redis failed")
)

type (
	// A RedisConf is a redis config.
	RedisConf struct {
		Host string `yaml:"Host"`
		Type string `json:",default=node,options=node|cluster" yaml:"Type"`
		Pass string `json:",optional" yaml:"Pass"`
		Tls  bool   `json:",optional" yaml:"Tls"`
	}

	// A RedisKeyConf is a redis config with key.
	RedisKeyConf struct {
		RedisConf
		Key string `json:",optional"`
	}
)

// Validate validates the RedisConf.
func (rc RedisConf) Validate() error {
	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}

	if len(rc.Type) == 0 {
		return ErrEmptyType
	}

	return nil
}

// Validate validates the RedisKeyConf.
func (rkc RedisKeyConf) Validate() error {
	if err := rkc.RedisConf.Validate(); err != nil {
		return err
	}

	if len(rkc.Key) == 0 {
		return ErrEmptyKey
	}

	return nil
}
