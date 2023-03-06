package kernel

import (
	"context"
	"github.com/flamego/flamego"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net/http"
	"pinnacle-primary-be/config"
	"pinnacle-primary-be/core/store/mysql"
	"pinnacle-primary-be/core/store/rds"
)

type (
	Engine struct {
		MainMysql  *mysql.Orm
		MainCache  *rds.Redis
		Fg         *flamego.Flame
		Grpc       *grpc.Server
		Conn       *grpc.ClientConn
		Mux        *runtime.ServeMux
		HttpServer *http.Server

		Ctx    context.Context
		Cancel context.CancelFunc

		ConfigListener []func(*config.GlobalConfig)
	}
)

var Kernel *Engine
