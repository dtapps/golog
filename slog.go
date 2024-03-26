package golog

import (
	"context"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
)

type SLogFun func() *SLog

type sLogConfig struct {
	showLine               bool              // 显示代码行
	lumberjackConfig       lumberjack.Logger // 配置lumberjack
	lumberjackConfigStatus bool
}

type SLog struct {
	option      sLogConfig
	jsonHandler *slog.JSONHandler
	logger      *slog.Logger
}

// NewSlog 创建
func NewSlog(opts ...SLogOption) *SLog {
	sl := &SLog{}
	for _, opt := range opts {
		opt(sl)
	}
	sl.start()
	return sl
}

func (sl *SLog) start() {

	opts := slog.HandlerOptions{
		AddSource: sl.option.showLine,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(gotime.DateTimeFormat))
				//return slog.Attr{}
			}
			return a
		},
	}

	// json格式输出
	var mw io.Writer
	if sl.option.lumberjackConfigStatus {
		// 同时控制台和文件输出日志
		mw = io.MultiWriter(os.Stdout, &sl.option.lumberjackConfig)
	} else {
		// 只在文件输出日志
		mw = io.MultiWriter(os.Stdout)
	}

	// 控制台输出
	sl.jsonHandler = slog.NewJSONHandler(mw, &opts)

	sl.logger = slog.New(sl.jsonHandler)

}

// WithLogger 跟踪编号
func (sl *SLog) WithLogger() *slog.Logger {
	logger := slog.New(sl.jsonHandler)
	return logger
}

// WithTraceId 跟踪编号
func (sl *SLog) WithTraceId(ctx context.Context) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String("trace_id", gotrace_id.GetTraceIdContext(ctx)),
		//slog.String("go_os", runtime.GOOS),
		//slog.String("go_arch", runtime.GOARCH),
		//slog.String("go_version", runtime.Version()),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceID 跟踪编号
func (sl *SLog) WithTraceID(ctx context.Context) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String("trace_id", gotrace_id.GetTraceIdContext(ctx)),
		//slog.String("go_os", runtime.GOOS),
		//slog.String("go_arch", runtime.GOARCH),
		//slog.String("go_version", runtime.Version()),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceIdStr 跟踪编号
func (sl *SLog) WithTraceIdStr(traceId string) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String("trace_id", traceId),
		//slog.String("go_os", runtime.GOOS),
		//slog.String("go_arch", runtime.GOARCH),
		//slog.String("go_version", runtime.Version()),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceIDStr 跟踪编号
func (sl *SLog) WithTraceIDStr(traceID string) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String("trace_id", traceID),
		//slog.String("go_os", runtime.GOOS),
		//slog.String("go_arch", runtime.GOARCH),
		//slog.String("go_version", runtime.Version()),
	})
	logger := slog.New(jsonHandler)
	return logger
}
