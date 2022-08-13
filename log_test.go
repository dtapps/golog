package golog

import (
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	g := NewGoLog(&ConfigGoLog{
		LogPath:      "./",
		LogName:      "all.log",
		LogLevel:     "debug",
		MaxSize:      2,
		MaxBackups:   30,
		MaxAge:       0,
		Compress:     false,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: true,
	})
	log.Println(g.Logger)
	g.Logger.Debug("debug 日志")
	g.Logger.Sugar().Debug("debug 日志")
	g.Logger.Info("info 日志")
	g.Logger.Sugar().Info("info 日志")
	g.Logger.Warn("warning 日志")
	g.Logger.Sugar().Warn("warning 日志")
	g.Logger.Error("error 日志")
	g.Logger.Sugar().Error("error 日志")
	log.Println(g.Logger)
}
