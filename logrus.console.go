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

	// 设置文件名和方法信息
	lrc.logger.SetReportCaller(true)

	lrc.logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: gotime.DateTimeFormat,
	})

	return lrc
}

// WithTraceId 跟踪编号
func (lrc *LogRusConsole) WithTraceId(ctx context.Context) *logrus.Entry {
	return lrc.logger.WithField("trace_id", gotrace_id.GetTraceIdContext(ctx))
}
