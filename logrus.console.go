package golog

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
)

type LogRusConsole struct {
	logger *logrus.Logger // 日志服务
	entry  *logrus.Entry  // 日志
}

// NewLogRusConsole 初始化控制台
func NewLogRusConsole() *LogRusConsole {

	lrc := &LogRusConsole{}
	lrc.logger = logrus.New()

	lrc.logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: gotime.DateTimeFormat,
	})

	return lrc
}

// 跟踪编号
func (lrc *LogRusConsole) withTraceId(ctx context.Context) {
	traceId := gotrace_id.GetTraceIdContext(ctx)
	if traceId == "" {
		lrc.entry = lrc.logger.WithFields(logrus.Fields{})
	} else {
		lrc.entry = lrc.logger.WithField("trace_id", gotrace_id.GetTraceIdContext(ctx))
	}
}

// Print 打印
func (lrc *LogRusConsole) Print(ctx context.Context, args ...interface{}) {
	lrc.withTraceId(ctx)
	lrc.entry.Print(args...)
}

// Printf 打印
func (lrc *LogRusConsole) Printf(ctx context.Context, format string, args ...interface{}) {
	lrc.withTraceId(ctx)
	lrc.entry.Printf(format, args...)
}

// Println 打印
func (lrc *LogRusConsole) Println(ctx context.Context, args ...interface{}) {
	lrc.withTraceId(ctx)
	lrc.entry.Println(args...)
}
