package applog

import (
	"fmt"
	"os"

	takin_adapter "github.com/YangZhaoWeblog/GoldenTakin/takin_log/adapter"
	takin_log_outpter "github.com/YangZhaoWeblog/GoldenTakin/takin_log/outputer"

	"github.com/YangZhaoWeblog/GoldenTakin/takin_log"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet 提供依赖注入组件
var ProviderSet = wire.NewSet(NewAppLogger, InitGlobalLogger)

// NewAppLogger 创建TakinLogger实例
func NewAppLogger() (*takin_log.TakinLogger, func()) {
	// 使用开箱即用方式配置文件日志
	opts := takin_log.AppLoggerOptions{
		Component: "test-component",
		AppName:   "test-app-builtin",
		MinLevel:  takin_log.DebugLevel,
		// 直接使用FileLogOption，无需手动创建fileOutput
		FileLogOption: &takin_log_outpter.FileLogOption{
			FilePath:   "./app_builtin.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
			LocalTime:  true,
		},
	}

	applogger := takin_log.NewAppLogger(opts)

	cleanUp := func() {
		err := applogger.Close()
		if err != nil {
			// 使用标准错误输出，避免在日志关闭时还使用日志
			fmt.Fprintf(os.Stderr, "failed to close applogger: %v\n", err)
		}
	}
	return applogger, cleanUp
}

// 初始化全局日志器
func InitGlobalLogger(takinLogger *takin_log.TakinLogger) log.Logger {
	// 创建适配器并设置为全局日志器
	adapter := takin_adapter.NewKratosAdapter(takinLogger)
	log.SetLogger(adapter)
	return adapter
}
