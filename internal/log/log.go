package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gopkg.in/natefinch/lumberjack.v2"
	"qqShuiHu/internal/conf"
)

// ProviderSet is a provider set for wire
var ProviderSet = wire.NewSet(NewLogger, NewLogHelper)

// Logger wraps slog.Logger
type Logger struct {
	l *slog.Logger
}

// NewLogger creates a new Logger instance
func NewLogger(c *conf.Log) log.Logger {
	// 确保日志目录存在
	logDir := c.Dir // 从配置中获取日志目录，如果没有配置则使用默认值
	if logDir == "" {
		logDir = "logs" // 默认日志目录
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("create log directory failed: %v", err))
	}

	// 配置日志轮转
	rotateLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "app.log"),
		MaxSize:    100,  // 每个日志文件最大尺寸，单位是 MB
		MaxBackups: 30,   // 保留的旧日志文件最大数量
		MaxAge:     7,    // 保留的旧日志文件最大天数
		Compress:   true, // 是否压缩旧日志文件
	}

	// 创建多写器，同时写入文件和标准输出
	writers := []io.Writer{rotateLogger}
	if c.Stdout { // 可以通过配置控制是否同时输出到标准输出
		writers = append(writers, os.Stdout)
	}
	multiWriter := io.MultiWriter(writers...)

	// 设置日志级别
	logLevel := slog.LevelInfo
	if c.Level != "" {
		switch c.Level {
		case "debug":
			logLevel = slog.LevelDebug
		case "info":
			logLevel = slog.LevelInfo
		case "warn":
			logLevel = slog.LevelWarn
		case "error":
			logLevel = slog.LevelError
		}
	}

	// 创建 JSON 处理器
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 自定义时间格式
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(time.Now().Format("2006-01-02 15:04:05.000")),
				}
			}
			return a
		},
	})

	// 创建 slog.Logger
	slogger := slog.New(handler)

	// 设置全局默认 logger
	slog.SetDefault(slogger)

	return &Logger{
		l: slogger,
	}
}

func NewLogHelper(logger log.Logger) *log.Helper {
	return log.NewHelper(logger)
}

// Log implements log.Logger
func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	var slogLevel slog.Level
	switch level {
	case log.LevelDebug:
		slogLevel = slog.LevelDebug
	case log.LevelInfo:
		slogLevel = slog.LevelInfo
	case log.LevelWarn:
		slogLevel = slog.LevelWarn
	case log.LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	var msg = ""
	attrs := make([]slog.Attr, 0)
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = fmt.Sprint(keyvals[i])
		}
		if key == "msg" {
			msg = fmt.Sprint(keyvals[i+1])
			continue
		}
		attrs = append(attrs, slog.Any(key, keyvals[i+1]))
	}

	l.l.LogAttrs(context.Background(), slogLevel, msg, attrs...)
	return nil
}

// Helper returns a new log.Helper
func (l *Logger) Helper() *log.Helper {
	return log.NewHelper(l)
}
