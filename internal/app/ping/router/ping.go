package router

import (
	"errors"
	"github.com/flamego/flamego"
	"pinnacle-primary-be/core/middleware/response"
)

func AppPingInit(e *flamego.Flame) {
	e.Get("/ping/v1", func(r flamego.Render) {
		response.HTTPSuccess(r, map[string]any{
			"message": "pong",
		})
	})

	e.Get("/ping/v1/err", func(r flamego.Render) {
		response.HTTPFail(r, 500000, "test error", errors.New("this is err"))
	})
}
