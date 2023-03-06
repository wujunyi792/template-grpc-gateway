package server

import (
	"context"
	"fmt"
	sentryflame "github.com/asjdf/flamego-sentry"
	"github.com/flamego/cors"
	"github.com/flamego/flamego"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"pinnacle-primary-be/config"
	"pinnacle-primary-be/core/healthz"
	"pinnacle-primary-be/core/kernel"
	"pinnacle-primary-be/core/logx"
	"pinnacle-primary-be/core/middleware/gw"
	"pinnacle-primary-be/core/middleware/rpc"
	"pinnacle-primary-be/core/sentryx"
	"pinnacle-primary-be/core/store/mysql"
	"pinnacle-primary-be/core/store/rds"
	"pinnacle-primary-be/core/stringx"
	"pinnacle-primary-be/core/tracex"
	"pinnacle-primary-be/internal/app/appInitialize"
	"pinnacle-primary-be/pkg/ip"
	"strings"
	"syscall"
	"time"
)

var (
	configYml string
	engine    *kernel.Engine
	StartCmd  = &cobra.Command{
		Use:     "server",
		Short:   "Set Application config info",
		Example: "main server -c config/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			setUp()
			loadStore()
			loadApp()
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/config.yaml", "Start server with provided configuration file")
}

// 初始化配置和日志
func setUp() {
	// 初始化全局 ctx
	ctx, cancel := context.WithCancel(context.Background())

	// 初始化资源管理器
	engine = &kernel.Engine{Ctx: ctx, Cancel: cancel}
	kernel.Kernel = engine

	// 加载配置
	config.LoadConfig(configYml, func(globalConfig *config.GlobalConfig) {
		for _, listener := range engine.ConfigListener {
			listener(globalConfig)
		}
	})

	// 设置日志等级
	if config.GetConfig().MODE == "" || config.GetConfig().MODE == "debug" {
		logx.SetLevel(logx.DebugLevel)
	}

	// 初始化 sentry
	sentryx.NewSentry(config.GetConfig().Sentry)

	// 初始化 opentelemetry
	tracex.StartAgent(config.GetConfig().Trace)

	// 初始化 flamego
	flamego.SetEnv(flamego.EnvType(config.GetConfig().MODE))
	engine.Fg = flamego.New()
	engine.Fg.Use(flamego.Recovery(), gw.RequestLog(), flamego.Renderer(), cors.CORS(cors.Options{
		AllowCredentials: true,
	}))
	if config.GetConfig().Sentry.Available() {
		engine.Fg.Use(sentryflame.New(sentryflame.Options{Repanic: true})) // sentry
	}

	// 初始化 grpc服务端
	engine.Grpc = grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		otelgrpc.UnaryServerInterceptor(),
		grpcrecovery.UnaryServerInterceptor(),
		grpcctxtags.UnaryServerInterceptor(),
		grpcauth.UnaryServerInterceptor(rpc.AuthInterceptor),
		rpc.LoggerInterceptor,
	)))
	reflection.Register(engine.Grpc)

	// 初始化 gateway grpc 客户端
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%s", config.GetConfig().Port), opts...)
	if err != nil {
		logx.Errorw("gRPC fail to dial", logx.LogField{Key: "err", Value: err})
		os.Exit(1)
	}

	// 初始化 gateway
	mux := runtime.NewServeMux(
		runtime.WithHealthzEndpoint(grpc_health_v1.NewHealthClient(conn)), // 健康检查
		runtime.WithIncomingHeaderMatcher(gw.IncomeMatcher),
		runtime.WithOutgoingHeaderMatcher(gw.OutgoingMatcher),
		runtime.WithErrorHandler(gw.GrpcGatewayError),            // 错误封装
		runtime.WithForwardResponseOption(gw.GrpcGatewaySuccess), // success 响应封装
		runtime.WithMarshalerOption("*", &gw.CustomMarshaller{}), // 为了实现将响应封装在固定格式json.data中，hack一下，在 ForwardResponseOption 中实现
	)
	engine.Mux = mux
	engine.Conn = conn
}

// 存储介质连接
func loadStore() {
	engine.MainMysql = mysql.MustNewMysqlOrm(config.GetConfig().MainMysql)
	engine.MainCache = rds.MustNewRedis(config.GetConfig().MainCache)
}

// 加载应用，包含多个生命周期
func loadApp() {
	apps := appInitialize.GetApps()
	for _, app := range apps {
		_err := app.PreInit(engine)
		if _err != nil {
			logx.Errorw("failed to pre init app", logx.LogField{Key: "error", Value: _err})
			os.Exit(1)
		}
	}
	for _, app := range apps {
		_err := app.Init(engine)
		if _err != nil {
			logx.Errorw("failed to init app", logx.LogField{Key: "error", Value: _err})
			os.Exit(1)
		}
	}
	for _, app := range apps {
		_err := app.PostInit(engine)
		if _err != nil {
			logx.Errorw("failed to post init app", logx.LogField{Key: "error", Value: _err})
			os.Exit(1)
		}
	}
	for _, app := range apps {
		_err := app.Load(engine)
		if _err != nil {
			logx.Errorw("failed to load app", logx.LogField{Key: "error", Value: _err})
			os.Exit(1)
		}
	}
	for _, app := range apps {
		_err := app.Start(engine)
		if _err != nil {
			logx.Errorw("failed to start app", logx.LogField{Key: "error", Value: _err})
			os.Exit(1)
		}
	}

	// 设置/grpc路由 将gw嵌入到flamego中，flamego 为入口网关，含 /grpc 前缀的请求转发到 grpc-gateway 处理
	engine.Fg.Any("/grpc/{**}", func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.Replace(r.RequestURI, "/grpc", "", 1)
		r.URL.Path = strings.Replace(r.URL.Path, "/grpc", "", 1)
		engine.Mux.ServeHTTP(w, r)
	})

}

// 启动服务
func run() {
	port := config.GetConfig().Port
	// 开启 tcp 监听
	conn, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logx.Errorw("failed to listen", logx.LogField{Key: "error", Value: err})
	}

	// 分流
	tcpMux := cmux.New(conn)
	grpcL := tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := tcpMux.Match(cmux.HTTP1Fast())
	go func() {
		// 在 flamego 外再包一层 otelhttp 用于链路追踪注入
		engine.HttpServer = &http.Server{
			Handler: otelhttp.NewHandler(engine.Fg, "gateway", otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			})),
		}
		if _err := engine.HttpServer.Serve(httpL); _err != nil && _err != http.ErrServerClosed {
			logx.Infow("failed to start to listen and serve http", logx.LogField{Key: "error", Value: _err})
		}
	}()
	go func() {
		if _err := engine.Grpc.Serve(grpcL); _err != nil {
			logx.Infow("failed to start to listen and serve grpc", logx.LogField{Key: "error", Value: _err})
		}
	}()

	go func() {
		logx.Info("mux listen starting...")
		if _err := tcpMux.Serve(); _err != nil {
			logx.Errorw("failed to serve mux", logx.LogField{Key: "error", Value: _err})
		}
	}()

	println(stringx.Green("Server run at:"))
	println(fmt.Sprintf("-  Local:   http://localhost:%s", port))
	for _, host := range ip.GetLocalHost() {
		println(fmt.Sprintf("-  Network: http://%s:%s", host, port))
	}
	// 健康检查设置为可接受服务
	healthz.Health.Set(true)

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 健康检查设置为不可接受服务
	healthz.Health.Set(false)

	println(stringx.Blue("Shutting down server..."))
	tracex.StopAgent()

	ctx, cancel := context.WithTimeout(engine.Ctx, 5*time.Second)
	defer engine.Cancel()
	defer cancel()

	if err := engine.HttpServer.Shutdown(ctx); err != nil {
		println(stringx.Yellow("Server forced to shutdown: " + err.Error()))
	}

	println(stringx.Green("Server exiting Correctly"))
}
