package golog

import (
	"context"
	"testing"
)

func TestNewSlog(t *testing.T) {
	sl := NewSlogConsole(context.Background(), &SLogConsoleConfig{
		ShowLine: true,
	})
	sl.WithLogger().Info("测试日志", "名称1", "内容1", "名称2", "内容2")
}
