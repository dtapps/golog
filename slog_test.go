package golog

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"testing"
)

func TestNewSlog(t *testing.T) {
	sl := NewSlog(
		WithSLogLumberjack(lumberjack.Logger{
			Filename:   "./test.log",
			MaxSize:    1,
			MaxAge:     3,
			MaxBackups: 4,
			LocalTime:  true,
			Compress:   true,
		}),
		WithSLogShowLine(),
	)
	sl.WithTraceIDStr("22").Info("测试日志", "名称1", "内容1", "名称2", "内容2")
}
