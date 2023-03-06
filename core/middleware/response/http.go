package response

import (
	"github.com/flamego/flamego"
)

type JsonResponse struct {
	Code    int32  `json:"code"`
	Error   any    `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func http(r flamego.Render, code int32, msg string, data any, err any) {
	r.JSON(int(code/1000), &JsonResponse{
		Code:    code,
		Message: msg,
		Data:    data,
		Error:   err,
	})
}

// HTTPSuccess 成功返回
func HTTPSuccess(r flamego.Render, data any) {
	http(r, 200000, "success", data, nil)
}

func HTTPFail(r flamego.Render, code int, msg string, err ...any) {
	for i, e := range err {
		if v, ok := e.(error); ok {
			err[i] = v.Error()
		}
	}
	http(r, int32(code), msg, nil, err)
}

func UnAuthorization(r flamego.Render) {
	HTTPFail(r, 401000, "登录过期失效，请重新登陆", nil)
}

func InValidParam(r flamego.Render, err ...any) {
	HTTPFail(r, 40200, "请求校验失败", err...)
}

func ServiceErr(r flamego.Render, err ...any) {
	HTTPFail(r, 500000, "内部异常", err...)
}
