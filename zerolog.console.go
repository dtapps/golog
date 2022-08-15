package golog

import (
	"context"
	"github.com/rs/zerolog"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
)

type ZeroLogConsole struct {
	logger zerolog.Logger        // 日志服务
	w      zerolog.ConsoleWriter // 日志
}

// NewZeroLogConsole 初始化控制台
func NewZeroLogConsole() *ZeroLogConsole {

	zlc := &ZeroLogConsole{}

	zlc.w = zerolog.NewConsoleWriter()
	zlc.w.TimeFormat = gotime.DateTimeFormat

	return zlc
}

// 跟踪编号
func (zlc *ZeroLogConsole) withTraceId(ctx context.Context) {
	traceId := gotrace_id.GetTraceIdContext(ctx)
	if traceId == "" {
		zlc.logger = zerolog.New(zlc.w).With().Timestamp().Logger()
	} else {
		zlc.logger = zerolog.New(zlc.w).With().Str("trace_id", gotrace_id.GetTraceIdContext(ctx)).Timestamp().Logger()
	}
}

// Print 打印
func (zlc *ZeroLogConsole) Print(ctx context.Context, v ...interface{}) {
	zlc.withTraceId(ctx)
	zlc.logger.Print(v...)
}

// Printf 打印
func (zlc *ZeroLogConsole) Printf(ctx context.Context, format string, v ...interface{}) {
	zlc.withTraceId(ctx)
	zlc.logger.Printf(format, v...)
}
