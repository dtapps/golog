package golog

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"io"
	"log"
	"os"
	"path"
)

type LogRusLogConfig struct {
	LogPath      string // 日志文件路径
	LogInConsole bool   // 是否同时输出到控制台
}

type LogRusLog struct {
	config *LogRusLogConfig // 配置
	logger *logrus.Logger   // 日志服务
	entry  *logrus.Entry    // 日志
	level  logrus.Level     // 日志等级
}

func NewLogRusLog(config *LogRusLogConfig) *LogRusLog {

	lr := &LogRusLog{config: config}
	lr.logger = logrus.New()

	// 设置为json格式
	lr.logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: gotime.DateTimeFormat,
	})

	// 需要保存到文件
	if lr.config.LogPath != "" {
		lr.logger.SetReportCaller(true)
	}

	return lr
}

// 跟踪编号
func (lr *LogRusLog) withTraceId(ctx context.Context) *LogRusLog {
	lr.entry = lr.logger.WithField("trace_id", gotrace_id.GetTraceIdContext(ctx))
	return lr
}

// Print 打印
func (lr *LogRusLog) Print(ctx context.Context, args ...interface{}) {
	lr.withTraceId(ctx)
	lr.entry.Print(args...)
}

// Printf 打印
func (lr *LogRusLog) Printf(ctx context.Context, format string, args ...interface{}) {
	lr.withTraceId(ctx)
	lr.entry.Printf(format, args...)
}

// Println 打印
func (lr *LogRusLog) Println(ctx context.Context, args ...interface{}) {
	lr.withTraceId(ctx)
	lr.entry.Println(args...)
}

// Panic 记录日志，然后panic
func (lr *LogRusLog) Panic(ctx context.Context, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Panic(args...)
}

// Panicf 记录日志，然后panic
func (lr *LogRusLog) Panicf(ctx context.Context, format string, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Panicf(format, args...)
}

// Fatal 有致命性错误，导致程序崩溃，记录日志，然后退出
func (lr *LogRusLog) Fatal(ctx context.Context, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Fatal(args...)
}

// Fatalf 有致命性错误，导致程序崩溃，记录日志，然后退出
func (lr *LogRusLog) Fatalf(ctx context.Context, format string, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Fatalf(format, args...)
}

// Error 错误日志
func (lr *LogRusLog) Error(ctx context.Context, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Error(args...)
}

// Errorf 错误日志
func (lr *LogRusLog) Errorf(ctx context.Context, format string, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Errorf(format, args...)
}

// Warn 警告日志
func (lr *LogRusLog) Warn(ctx context.Context, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Warn(args...)
}

// Warnf 警告日志
func (lr *LogRusLog) Warnf(ctx context.Context, format string, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Warnf(format, args...)
}

// Info 核心流程日志
func (lr *LogRusLog) Info(ctx context.Context, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Info(args...)
}

// Infof 核心流程日志
func (lr *LogRusLog) Infof(ctx context.Context, format string, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Infof(format, args...)
}

// Debug debug日志（调试日志）
func (lr *LogRusLog) Debug(ctx context.Context, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Debug(args...)
}

// Debugf debug日志（调试日志）
func (lr *LogRusLog) Debugf(ctx context.Context, format string, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Debugf(format, args...)
}

// Trace 粒度超细的，一般情况下我们使用不上
func (lr *LogRusLog) Trace(ctx context.Context, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Debug(args...)
}

// Tracef 粒度超细的，一般情况下我们使用不上
func (lr *LogRusLog) Tracef(ctx context.Context, format string, args ...interface{}) {
	lr.setOutPutFile()
	lr.withTraceId(ctx)
	lr.entry.Tracef(format, args...)
}

// https://www.fushengwushi.com/archives/1397
// https://blog.csdn.net/oscarun/article/details/114295955
// https://juejin.cn/post/6974353191974469668
func (lr *LogRusLog) setOutPutFile() {

	// 是否保存到文件
	if lr.config.LogPath == "" {
		return
	}

	// 判断文件夹
	if _, err := os.Stat(lr.config.LogPath); os.IsNotExist(err) {
		err = os.MkdirAll(lr.config.LogPath, 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: %s", lr.config.LogPath, err))
		}
	}

	// 日志名
	fileName := path.Join(lr.config.LogPath, "logrus."+gotime.Current().SetFormat(gotime.ShortDateFormat)+".log")

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	// 是否同时输出到控制台
	if lr.config.LogInConsole {
		lr.logger.SetOutput(os.Stdout)
		writers := []io.Writer{
			file,
			os.Stdout,
		}
		fileAndStdoutWriter := io.MultiWriter(writers...)
		if err == nil {
			lr.logger.SetOutput(fileAndStdoutWriter)
		} else {
			log.Printf("无法记录到文件 %s\n", fileName)
		}
	} else {
		lr.logger.SetOutput(file)
	}

	return
}
