package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/YangZhaoWeblog/GoldenTakin/takin_log"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

type logArgs struct {
	Kind      string
	Component string
	Operation string
	Args      string
	Code      int
	Reason    string
	Stack     string
	Latency   float64
}

func (l logArgs) toKV() []interface{} {
	return []interface{}{
		"kind", l.Kind,
		"component", l.Component,
		"operation", l.Operation,
		"args", l.Args,
		"code", l.Code,
		"reason", l.Reason,
		"stack", l.Stack,
		"latency", l.Latency,
	}
}

// ServerLog is an server logging middleware.
func ServerLog(appLogger *takin_log.TakinLogger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 1. 初始化日志参数
			args := newLogArgs("client", req)

			// 2. 获取请求上下文信息
			if info, ok := transport.FromServerContext(ctx); ok {
				args.Kind = info.Kind().String()
				args.Operation = info.Operation()
			}

			// 3. 执行请求并记录日志
			startTime := time.Now()
			defer func() {
				logRequestResult(ctx, appLogger, args, err, startTime)
			}()
			reply, err = handler(ctx, req)
			return
		}
	}
}

// ClientLog is a client logging middleware.
func ClientLog(appLogger *takin_log.TakinLogger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 1. 初始化日志参数
			args := newLogArgs("client", req)

			// 2. 获取请求上下文信息
			if info, ok := transport.FromClientContext(ctx); ok {
				args.Kind = info.Kind().String()
				args.Operation = info.Operation()
			}

			// 3. 执行请求并记录日志
			startTime := time.Now()
			defer func() {
				logRequestResult(ctx, appLogger, args, err, startTime)
			}()
			reply, err = handler(ctx, req)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

func newLogArgs(kind string, req interface{}) *logArgs {
	return &logArgs{
		Kind:      kind,
		Component: "",
		Operation: "",
		Args:      extractArgs(req),
		Code:      0,
		Reason:    "",
		Stack:     "",
	}
}

func logRequestResult(ctx context.Context, logger *takin_log.TakinLogger, args *logArgs, err error, startTime time.Time) {
	var msg string
	if se := errors.FromError(err); se != nil {
		args.Code = int(se.Code)
		args.Reason = se.Reason
		msg = se.Message
	}

	args.Latency = time.Since(startTime).Seconds()
	// 记录错误日志还是正常日志
	if err != nil {
		args.Stack = fmt.Sprintf("%+v", err)
		logger.ErrorContext(ctx, msg, args.toKV()...)
		return
	}
	logger.InfoContext(ctx, msg, args.toKV()...)
}
