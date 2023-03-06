package web

import (
	"github.com/asaskevich/govalidator"
	"github.com/flamego/flamego"
	"pinnacle-primary-be/core/jsonx"
	"pinnacle-primary-be/core/middleware/response"
	"pinnacle-primary-be/internal/websocket"
	"pinnacle-primary-be/pkg/utils/page"
	"reflect"
)

func InjectWebsocket(key ...string) flamego.Handler {
	if len(key) > 1 {
		panic("InjectWebsocket only accept at most one key")
	}
	if len(key) == 1 {
		return func(c flamego.Context) {
			c.Map(websocket.GetSocketManager(key[0]))
		}
	}
	return func(c flamego.Context) {
		c.Map(websocket.GetSocketManager("*"))
	}
}

func InjectJson[T any]() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req T
		body, err := c.Request().Body().Bytes()
		if err = jsonx.Unmarshal(body, &req); err != nil {
			response.InValidParam(r, err)
			return
		}
		_, err = govalidator.ValidateStruct(&req)
		if err != nil {
			response.InValidParam(r, err)
			return
		}
		c.Map(req)
	}
}

func InjectQuery[T any]() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req T
		t := reflect.TypeOf(req)
		v := reflect.ValueOf(&req).Elem()
		var tag reflect.StructTag
		for i := 0; i < v.NumField(); i++ {
			tag = t.Field(i).Tag
			if value, ok := tag.Lookup("query"); ok && t.Field(i).IsExported() {
				switch v.Field(i).Kind() {
				case reflect.String:
					v.Field(i).SetString(c.Query(value))
				case reflect.Int, reflect.Int64:
					v.Field(i).SetInt(c.QueryInt64(value))
				case reflect.Bool:
					v.Field(i).SetBool(c.QueryBool(value))
				}
			}
		}

		_, err := govalidator.ValidateStruct(&req)
		if err != nil {
			response.InValidParam(r, err)
			return
		}
		c.Map(req)
	}
}

func InjectParam[T any]() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req T
		t := reflect.TypeOf(req)
		v := reflect.ValueOf(&req).Elem()
		var tag reflect.StructTag
		for i := 0; i < v.NumField(); i++ {
			tag = t.Field(i).Tag
			if value, ok := tag.Lookup("param"); ok && t.Field(i).IsExported() {
				switch v.Field(i).Kind() {
				case reflect.String:
					v.Field(i).SetString(c.Param(value))
				case reflect.Int, reflect.Int64:
					v.Field(i).SetInt(c.ParamInt64(value))
				case reflect.Bool:
					panic("bool type not support")
				}
			}
		}

		_, err := govalidator.ValidateStruct(&req)
		if err != nil {
			response.InValidParam(r, err)
			return
		}
		c.Map(req)
	}
}

func InjectPaginate() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req page.Paginate
		req.Current = c.QueryInt("current")
		req.PageSize = c.QueryInt("pageSize")
		c.Map(req)
	}
}
