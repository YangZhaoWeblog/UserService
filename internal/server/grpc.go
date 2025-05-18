package server

import (
	v1 "github.com/YangZhaoWeblog/UserService/api/helloworld/v1"
	userv1 "github.com/YangZhaoWeblog/UserService/api/user/v1"
	"github.com/YangZhaoWeblog/UserService/internal/conf"
	"github.com/YangZhaoWeblog/UserService/internal/observability"
	"github.com/YangZhaoWeblog/UserService/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService,
	user *service.UserService,
	metricsData *observability.MetricsData,
	tracer *sdktrace.TracerProvider,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			//recovery.Recovery(), //自动捕获 panic 确保线上服务不崩溃，测试环境应当尽可能让崩溃
			tracing.Server(),
			metrics.Server(
				metrics.WithSeconds(metricsData.Seconds),
				metrics.WithRequests(metricsData.Requests),
			),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	v1.RegisterGreeterServer(srv, greeter)
	userv1.RegisterUserServer(srv, user)
	return srv
}
