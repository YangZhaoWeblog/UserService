package observability

import (
	"github.com/YangZhaoWeblog/UserService/internal/conf"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// MetricsData 保存所有指标
type MetricsData struct {
	Seconds  metric.Float64Histogram
	Requests metric.Int64Counter
}

// 为什么高版本Kratos要用OpenTelemetry？
// 1. 统一标准: OpenTelemetry统一了指标、跟踪和日志的标准，是CNCF的核心项目, 是未来的标准
// 2. OpenTelemetry 可以对接多个后端，如 Jaeger, Prometheus, Zipkin, Datadog 等，而 Prometheus 只是其中一种
// 所以使用 OpenTelemetry 与 Grafana 并不冲突, 数据流向如下:
// 应用 → OpenTelemetry SDK → Prometheus导出器 → Prometheus服务器 → Grafana面板

// NewMetrics 提供指标信息，在此可以定义要捕获的各种指标信息， 用于wire注入
func NewMetrics(app *conf.App) (*MetricsData, error) {
	// 使用 OpenTelemetry 作为指标收集框架，但数据会被导出到 Prometheus 格式
	// 注意：配置好这个基础框架后，许多系统指标会被自动收集，不需要手动定义

	// 1. 创建 Prometheus 导出器
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	// 2. 创建 OpenTelemetry 指标提供者
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter := provider.Meter(app.ServiceName)

	// 以下是Kratos框架提供的标准指标，而非仅有的两个指标
	// OpenTelemetry 会自动收集更多系统和运行时指标，不需要手动添加

	// 核心指标1: 请求计数器，此为Kratos框架标准指标
	// 框架会自动收集并添加kind(http/grpc)、operation(路径/方法)、code(状态码)、reason(错误原因)标签
	requests, err := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName) //名称为: server_requests_code_total
	if err != nil {
		return nil, err
	}

	// 核心指标2: 请求处理时间直方图，此为Kratos框架标准指标
	// 框架会自动计算各分位数(P50/P95/P99)，以监控服务性能
	seconds, err := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		return nil, err
	}

	// 通过上述配置，已经启用了完整的指标收集系统
	// 除了这两个核心HTTP/gRPC指标外，还会自动收集Go运行时指标(GC、内存、goroutine等)
	// 其他添加业务指标，可以使用meter创建额外的计数器、仪表盘或直方图

	return &MetricsData{
		Seconds:  seconds,
		Requests: requests,
	}, nil
}
