package rds

import (
	"crypto/tls"
	red "github.com/redis/go-redis/v9"
	"io"
	"pinnacle-primary-be/core/syncx"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var clientManager = syncx.NewResourceManager()

func getClient(r *Redis) (*red.Client, error) {
	val, err := clientManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Password:     r.Pass,
			DB:           defaultDatabase,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
