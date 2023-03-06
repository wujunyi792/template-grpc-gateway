package web

import (
	"github.com/flamego/flamego"
	"pinnacle-primary-be/core/middleware/response"
	"pinnacle-primary-be/pkg/jwt"
)

func Authorization(c flamego.Context, r flamego.Render) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		response.UnAuthorization(r)
		return
	}
	entry, err := jwt.ParseToken(token)
	if err != nil {
		response.UnAuthorization(r)
		return
	}
	c.Map(entry.Info)
}
