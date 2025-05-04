// cmd/server/wire.go
package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/YangZhaoWeblog/GoldenTakin/takin_log"
	"github.com/YangZhaoWeblog/UserService/internal/conf"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0" // 使用语义约定包
)

// NewTracerProvider 创建并配置 OTel Tracer Provider
func NewTracerProvider(cfg *conf.Trace, cfgApp *conf.App, logHelper *takin_log.TakinLogger) (*sdktrace.TracerProvider, func(), error) {
	if cfg == nil || cfg.Endpoint == "" {
		return nil, func() {}, nil
	}

	ctx := context.Background()

	// 1. 创建 OTLP gRPC Exporter
	// 注意：生产环境应配置 TLS 和可能的认证
	// insecure.NewCredentials() 仅用于本地或受信任的网络
	// 你可能需要根据实际情况添加 grpc.WithTransportCredentials(...)
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(),             // 生产环境需要移除或替换为 TLS 配置
		otlptracegrpc.WithTimeout(5*time.Second), // 添加超时设置
		// otlptracegrpc.WithDialOption(grpc.WithBlock()), // 可选: 启动时阻塞直到连接建立
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create OTLP trace exporter failed: %w", err)
	}

	// 2. 配置资源属性 (Resource Attributes)
	// 这些属性会附加到所有由此 Provider 生成的 Span 上
	resAttrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(cfgApp.ServiceName), // service.name 是必须的
		// semconv.ServiceVersionKey.String(cfg.ServiceVersion), // 如果配置了版本
		// attribute.String("environment", cfg.Environment),   // 如果配置了环境
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(resAttrs...),
		resource.WithFromEnv(),      // 尝试从 OTEL_RESOURCE_ATTRIBUTES 环境变量读取
		resource.WithProcess(),      // 添加进程信息 (PID, 可执行文件名等)
		resource.WithTelemetrySDK(), // 添加 OTel SDK 信息 (名称, 版本, 语言)
		resource.WithHost(),         // 添加主机信息 (名称, 架构)
	)
	if err != nil {
		_ = exporter.Shutdown(ctx) // 忽略关闭错误
		return nil, nil, fmt.Errorf("create resource failed: %w", err)
	}

	// 3. 配置采样器 (Sampler)
	var sampler sdktrace.Sampler
	ratio := cfg.SampleRatio
	if ratio >= 1.0 {
		sampler = sdktrace.AlwaysSample() // 全部采样
	} else if ratio <= 0.0 {
		sampler = sdktrace.NeverSample() // 全部丢弃
	} else {
		// 基于父 Span 的采样状态，如果父 Span 被采样，则子 Span 也被采样
		// 对于根 Span (没有父 Span)，使用 TraceIDRatioBased 决定是否采样
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))
	}

	// 5. 创建 TracerProvider
	bsp := sdktrace.NewBatchSpanProcessor(exporter) // BatchSpanProcessor 将 Span 批量发送，提高性能
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp), // 使用 BatchSpanProcessor
		// sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter)), // SimpleSpanProcessor 用于调试，不推荐生产环境
	)

	// 6. 设置全局 TracerProvider 和 Propagator
	// 这使得应用中其他地方可以通过 otel.Tracer() 获取 Tracer
	// 并使得上下文可以在服务间传播 (HTTP Headers, gRPC Metadata)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // W3C Trace Context standard
		propagation.Baggage{},      // W3C Baggage standard
	))

	logHelper.Info("Tracer provider already setting", "service=", cfgApp.ServiceName, "endpoint=", cfg.Endpoint, "sample_ratio=", cfg.SampleRatio)

	// 返回 Provider 和一个清理函数，用于应用关闭时调用
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置关闭超时
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logHelper.Info("WARN: close tracer provider failed: %v", err)
		}
		logHelper.Info("Tracer provider already closed")
	}

	return tp, cleanup, nil
}
