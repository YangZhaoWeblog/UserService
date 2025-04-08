package applog

import (
	"fmt"
	"os"

	"github.com/YangZhaoWeblog/GoldenTakin/takin_log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewApplog)

// 测试使用内置FileLogOption的开箱即用方式
func NewApplog() (*takin_log.TakinLogger, func()) {
	// 使用开箱即用方式配置文件日志
	opts := takin_log.AppLoggerOptions{
		Component: "test-component",
		AppName:   "test-app-builtin",
		MinLevel:  takin_log.DebugLevel,
		// 直接使用FileLogOption，无需手动创建fileOutput
		FileLogOption: &takin_log.FileLogOption{
			FilePath:   "./app_builtin.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
			LocalTime:  true,
		},
	}

	applogger := takin_log.NewAppLoggerWithKratos(opts)

	cleanUp := func() {
		err := applogger.Close()
		if err != nil {
			// 使用标准错误输出，避免在日志关闭时还使用日志
			fmt.Fprintf(os.Stderr, "failed to close applogger: %v\n", err)
		}
	}
	return applogger, cleanUp
}
