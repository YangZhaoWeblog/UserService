package applog

import (
	"fmt"
	"github.com/YangZhaoWeblog/UserService/internal/conf"
	"os"

	takin_adapter "github.com/YangZhaoWeblog/GoldenTakin/takin_log/adapter"
	takin_log_outpter "github.com/YangZhaoWeblog/GoldenTakin/takin_log/outputer"

	"github.com/YangZhaoWeblog/GoldenTakin/takin_log"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAppLogger, InitGlobalLogger)

// NewAppLogger 创建TakinLogger实例
func NewAppLogger(confLog *conf.Log, confApp *conf.App) (*takin_log.TakinLogger, func()) {
	level, _ := takin_log.ParseLogLevel(confLog.Level)
	opts := takin_log.TakinLoggerOptions{
		Component: confApp.GetServiceName(),
		AppName:   confApp.GetAppName(),
		MinLevel:  level,
		FileLogOption: &takin_log_outpter.FileLogOption{
			FilePath:   confLog.Dir,
			MaxSize:    int(confLog.MaxSize),
			MaxBackups: int(confLog.MaxBackups),
			MaxAge:     int(confLog.MaxAge),
			Compress:   confLog.Compress,
		},
	}

	applogger := takin_log.NewTakinLogger(opts)

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
