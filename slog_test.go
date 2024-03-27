package golog

import (
	"context"
	"go.dtapp.net/gotrace_id"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
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
	sl.WithTraceIDStr("22").Info("测试链式日志", "名称1", "内容1", "名称2", "内容2")
}

func TestNewSlogSetDefault(t *testing.T) {
	NewSlog(
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
	slog.Info("测试默认日志", "名称1", "内容1", "名称2", "内容2")
	var ctx = gotrace_id.SetCustomTraceId(context.Background(), "23456")
	slog.InfoContext(ctx, "测试默认日志带上下文", "名称1", "内容1", "名称2", "内容2")
}

func TestNewSlogSetDefaultCtx(t *testing.T) {
	NewSlog(
		WithSLogLumberjack(lumberjack.Logger{
			Filename:   "./test.log",
			MaxSize:    1,
			MaxAge:     3,
			MaxBackups: 4,
			LocalTime:  true,
			Compress:   true,
		}),
		WithSLogShowLine(),
		WithSLogSetDefault(),
		WithSLogSetDefaultCtx(),
	)

	var ctx = gotrace_id.SetCustomTraceId(context.Background(), "23456")
	slog.InfoContext(ctx, "测试默认日志带上下文 InfoContext ", "名称1", "内容1", "名称2", "内容2")

}
