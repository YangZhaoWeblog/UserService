package server

import (
	"github.com/YangZhaoWeblog/GoldenTakin/takin_log"
	v1 "github.com/YangZhaoWeblog/UserService/api/helloworld/v1"
	userv1 "github.com/YangZhaoWeblog/UserService/api/user/v1"
	"github.com/YangZhaoWeblog/UserService/internal/conf"
	"github.com/YangZhaoWeblog/UserService/internal/observability"
	"github.com/YangZhaoWeblog/UserService/internal/server/middleware"
	"github.com/YangZhaoWeblog/UserService/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService,
	user *service.UserService,
	metricsData *observability.MetricsData,
	applogger *takin_log.TakinLogger,
	tracer *sdktrace.TracerProvider,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			//recovery.Recovery(), //自动捕获 panic 确保线上服务不崩溃，测试环境应当尽可能让崩溃
			tracing.Server(),
			metrics.Server( //指标中间件
				metrics.WithSeconds(metricsData.Seconds),
				metrics.WithRequests(metricsData.Requests),
			),
			middleware.ServerLog(applogger),
		),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	// Prometheus 定期访问 http://host.docker.internal:8010/metrics, 抓取收集的指标
	// 仅需再次 http server 暴露，grpc 的请求也会被自动捕获到
	srv.Handle("/metrics", promhttp.Handler())

	v1.RegisterGreeterHTTPServer(srv, greeter)
	userv1.RegisterUserHTTPServer(srv, user)
	return srv
}
