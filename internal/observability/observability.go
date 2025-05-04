package observability

import "github.com/google/wire"

//var ProviderSet = wire.NewSet(NewMetrics, InitGlobalTracer)

//var ProviderSet = wire.NewSet(NewAppLogger, InitGlobalLogger)

// 将 NewAppLogger、InitGlobalLogger、NewMetrics、NewTracerProvider, 全部串到一个 ProviderSet 中，并且合并它们的 cleanup
var ProviderSet = wire.NewSet(
	// 1. 日志, 被 kratos log 所依赖，所以无需被显式使用
	NewAppLogger,
	InitGlobalLogger,

	// 2. 指标（无 cleanup）
	NewMetrics,

	// 3. 追踪
	NewTracerProvider,
)
