package gw

import (
	"github.com/flamego/flamego"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"pinnacle-primary-be/core/logx"
	"time"
)

func RequestLog() flamego.Handler {
	return func(c flamego.Context, r *http.Request) {

		// 开始时间
		startTime := time.Now()

		otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(c.ResponseWriter().Header()))
		// 处理请求
		c.Next()

		logx.WithContext(r.Context()).Infow("request log", logx.Field("path", c.Request().RequestURI), logx.Field("method", c.Request().Method), logx.Field("ip", c.RemoteAddr()), logx.Field("status", c.ResponseWriter().Status()), logx.Field("duration", time.Now().Sub(startTime)))
	}
}
