package ping

import (
	"context"
	"errors"
	"pinnacle-primary-be/core/kernel"
	pingV1 "pinnacle-primary-be/gen/proto/ping/v1"
	"pinnacle-primary-be/internal/app/ping/dao"
	"pinnacle-primary-be/internal/app/ping/router"
	"pinnacle-primary-be/internal/app/ping/service/gw"
	"sync"
)

var (
	ErrEmptyDatabase = errors.New("database pointer is nil")
)

type (
	Ping struct {
		Name string
	}
)

func (p *Ping) Info() string {
	return p.Name
}

func (p *Ping) PreInit(engine *kernel.Engine) error {
	if engine.MainMysql == nil {
		return ErrEmptyDatabase
	}
	dao.Ping = engine.MainMysql
	return nil
}

func (p *Ping) Init(*kernel.Engine) error {
	_err := dao.AutoMigrate()
	if _err != nil {
		return _err
	}
	return nil
}

func (p *Ping) PostInit(*kernel.Engine) error {
	return nil
}

func (p *Ping) Load(engine *kernel.Engine) error {
	// 加载flamego api
	router.AppPingInit(engine.Fg)
	// 加载grpc gw
	pingV1.RegisterPingServiceServer(engine.Grpc, &gw.S{})
	_err := pingV1.RegisterPingServiceHandler(engine.Ctx, engine.Mux, engine.Conn)
	if _err != nil {
		return _err
	}
	return nil
}

func (p *Ping) Start(engine *kernel.Engine) error {
	return nil
}

func (p *Ping) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Ping) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {
		if engine.MainMysql == nil {
			return ErrEmptyDatabase
		}
		dao.Ping = engine.MainMysql
		return nil
	}
}
