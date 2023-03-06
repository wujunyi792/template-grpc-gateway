package app

import (
	"context"
	"pinnacle-primary-be/core/kernel"
	"sync"
)

type (
	Module interface {
		Info() string
		PreInit(*kernel.Engine) error
		Init(*kernel.Engine) error
		PostInit(*kernel.Engine) error
		Load(*kernel.Engine) error
		Start(*kernel.Engine) error
		Stop(wg *sync.WaitGroup, ctx context.Context) error

		OnConfigChange() func(*kernel.Engine) error
	}
)
