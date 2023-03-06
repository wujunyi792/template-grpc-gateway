package healthz

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"pinnacle-primary-be/core/healthz"
	"pinnacle-primary-be/core/kernel"
	"sync"
)

type (
	Healthz struct {
		Name string
	}
)

func (p *Healthz) Info() string {
	return p.Name
}

func (p *Healthz) PreInit(engine *kernel.Engine) error {
	return nil
}

func (p *Healthz) Init(*kernel.Engine) error {
	return nil
}

func (p *Healthz) PostInit(*kernel.Engine) error {
	return nil
}

func (p *Healthz) Load(engine *kernel.Engine) error {
	grpc_health_v1.RegisterHealthServer(engine.Grpc, &healthz.S{})
	return nil
}

func (p *Healthz) Start(engine *kernel.Engine) error {
	return nil
}

func (p *Healthz) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Healthz) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {
		return nil
	}
}
