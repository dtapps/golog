package golog

import (
	"context"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"testing"
)

// 参考 https://www.cnblogs.com/cheyunhua/p/18049634

func TestNewSlog1(t *testing.T) {
	sl := NewSlog(
		WithSLogLumberjack(&lumberjack.Logger{
			Filename:   "./test.log",
			MaxSize:    1,
			MaxAge:     3,
			MaxBackups: 4,
			LocalTime:  true,
			Compress:   true,
		}),
		WithSLogShowLine(),
	)
	sl.WithLogger().Info("测试链式日志", "名称1", "内容1", "名称2", "内容2")
}

func TestNewSlog2(t *testing.T) {
	NewSlog(
		WithSLogLumberjack(&lumberjack.Logger{
			Filename:   "./test.log",
			MaxSize:    1,
			MaxAge:     3,
			MaxBackups: 4,
			LocalTime:  true,
			Compress:   true,
		}),
		//WithSLogShowLine(),
		WithSLogSetDefault(),
		WithSLogSetDefaultCtx(),
		//WithSLogSetJSONFormat(),
	)
	slog.Info("测试默认日志", "名称1", "内容1", "名称2", "内容2")
	slog.InfoContext(context.TODO(), "测试默认日志带上下文", "名称1", "内容1", "名称2", "内容2")
}

func TestNewSlog3(t *testing.T) {
	NewSlog(
		WithSLogLumberjack(&lumberjack.Logger{
			Filename:   "./test.log",
			MaxSize:    1,
			MaxAge:     3,
			MaxBackups: 4,
			LocalTime:  true,
			Compress:   true,
		}),
		//WithSLogShowLine(),
		WithSLogSetDefault(),
		WithSLogSetDefaultCtx(),
		//WithSLogSetJSONFormat(),
	)

	slog.InfoContext(context.TODO(), "测试默认日志带上下文 InfoContext ", "名称1", "内容1", "名称2", "内容2")

}
