package golog

import (
	"context"
	"github.com/natefinch/lumberjack"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"log/slog"
	"os"
	"runtime"
)

type SLogFun func() *SLog

type sLogConfig struct {
	logPath    string // [File]日志文件路径
	logName    string // [File]日志文件名
	maxSize    int    // [File]单位为MB,默认为512MB
	maxBackups int    // [File]保留旧文件的最大个数
	maxAge     int    // [File]文件最多保存多少天 0=不删除
	showLine   bool   // [File、Console]显示代码行
}

type SLog struct {
	config      *sLogConfig
	jsonHandler *slog.JSONHandler
	textHandler *slog.TextHandler
	logger      *slog.Logger
}

type SLogFileConfig struct {
	LogPath    string // 日志文件路径
	LogName    string // 日志文件名
	MaxSize    int    // 单位为MB,默认为512MB
	MaxBackups int    // 保留旧文件的最大个数
	MaxAge     int    // 文件最多保存多少天 0=不删除
	ShowLine   bool   // 显示代码行
}

func NewSlogFile(ctx context.Context, config *SLogFileConfig) *SLog {

	sl := &SLog{
		config: &sLogConfig{
			logPath:    config.LogPath,
			logName:    config.LogName,
			maxSize:    config.MaxSize,
			maxBackups: config.MaxBackups,
			maxAge:     config.MaxAge,
			showLine:   config.ShowLine,
		},
	}

	opts := slog.HandlerOptions{
		AddSource: sl.config.showLine,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(gotime.DateTimeFormat))
				//return slog.Attr{}
			}
			return a
		},
	}

	lumberjackLogger := lumberjack.Logger{
		Filename:   sl.config.logPath + sl.config.logName, // ⽇志⽂件路径
		MaxSize:    sl.config.maxSize,                     // 单位为MB,默认为512MB
		MaxAge:     sl.config.maxAge,                      // 文件最多保存多少天
		MaxBackups: sl.config.maxBackups,                  // 保留旧文件的最大个数
	}

	// json格式输出
	sl.jsonHandler = slog.NewJSONHandler(&lumberjackLogger, &opts)
	sl.logger = slog.New(sl.jsonHandler)

	return sl
}

type SLogConsoleConfig struct {
	ShowLine bool // 显示代码行
}

func NewSlogConsole(ctx context.Context, config *SLogConsoleConfig) *SLog {

	sl := &SLog{
		config: &sLogConfig{
			showLine: config.ShowLine,
		},
	}

	opts := slog.HandlerOptions{
		AddSource: sl.config.showLine,
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
	sl.jsonHandler = slog.NewJSONHandler(os.Stdout, &opts)
	sl.logger = slog.New(sl.jsonHandler)

	return sl
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
		slog.String("go_os", runtime.GOOS),
		slog.String("go_arch", runtime.GOARCH),
		slog.String("go_version", runtime.Version()),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceID 跟踪编号
func (sl *SLog) WithTraceID(ctx context.Context) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String("trace_id", gotrace_id.GetTraceIdContext(ctx)),
		slog.String("go_os", runtime.GOOS),
		slog.String("go_arch", runtime.GOARCH),
		slog.String("go_version", runtime.Version()),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceIdStr 跟踪编号
func (sl *SLog) WithTraceIdStr(traceId string) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String("trace_id", traceId),
		slog.String("go_os", runtime.GOOS),
		slog.String("go_arch", runtime.GOARCH),
		slog.String("go_version", runtime.Version()),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceIDStr 跟踪编号
func (sl *SLog) WithTraceIDStr(traceID string) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String("trace_id", traceID),
		slog.String("go_os", runtime.GOOS),
		slog.String("go_arch", runtime.GOARCH),
		slog.String("go_version", runtime.Version()),
	})
	logger := slog.New(jsonHandler)
	return logger
}
