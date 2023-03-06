package router

import (
	"github.com/flamego/flamego"
)

func AppExampleInit(e *flamego.Flame) {
	e.Group("/api/v1/example", func() {
		ExampleGroup(e)
	})
}

func ExampleGroup(e *flamego.Flame) {}
