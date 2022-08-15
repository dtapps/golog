package golog

import (
	"context"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapLogConfig struct {
	LogPath      string // 日志文件路径
	LogName      string // 日志文件名
	LogLevel     string // 日志级别 debug/info/warn/error，debug输出：debug/info/warn/error日志。 info输出：info/warn/error日志。 warn输出：warn/error日志。 error输出：error日志。
	MaxSize      int    // 单个文件大小,MB
	MaxBackups   int    // 保存的文件个数
	MaxAge       int    // 保存的天数 0=不删除
	Compress     bool   // 压缩
	JsonFormat   bool   // 是否输出为json格式
	ShowLine     bool   // 显示代码行
	LogInConsole bool   // 是否同时输出到控制台
}

type ZapLog struct {
	config *ZapLogConfig
	logger *zap.Logger
}

func NewZapLog(config *ZapLogConfig) *ZapLog {

	zl := &ZapLog{config: config}

	// 设置日志级别
	var level zapcore.Level
	switch zl.config.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	var (
		syncer zapcore.WriteSyncer

		// 自定义时间输出格式
		customTimeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}

		// 自定义日志级别显示
		customLevelEncoder = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(level.CapitalString())
		}
	)

	// 定义日志切割配置
	hook := lumberjack.Logger{
		Filename:   zl.config.LogPath + zl.config.LogName, // 日志文件的位置
		MaxSize:    zl.config.MaxSize,                     // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: zl.config.MaxBackups,                  // 保留旧文件的最大个数
		Compress:   zl.config.Compress,                    // 是否压缩 disabled by default
	}
	if zl.config.MaxAge > 0 {
		hook.MaxAge = zl.config.MaxAge // days
	}

	// 判断是否控制台输出日志
	if zl.config.LogInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.AddSync(&hook)
	}

	// 定义zap配置信息
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,          // 自定义时间格式
		EncodeLevel:    customLevelEncoder,         // 小写编码器
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder

	// 判断是否json格式输出
	if zl.config.JsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)

	zl.logger = zap.New(core)

	// 判断是否显示代码行号
	if zl.config.ShowLine {
		zl.logger = zl.logger.WithOptions(zap.AddCaller())
	}

	return zl
}

// Panic 记录日志，然后panic
func (zl *ZapLog) Panic(ctx context.Context, args ...interface{}) {
	zl.logger.Sugar().Panic(args...)
}

// Panicf 记录日志，然后panic
func (zl *ZapLog) Panicf(ctx context.Context, template string, args ...interface{}) {
	zl.logger.Sugar().Panicf(template, args...)
}

// Fatal 有致命性错误，导致程序崩溃，记录日志，然后退出
func (zl *ZapLog) Fatal(ctx context.Context, args ...interface{}) {
	zl.logger.Sugar().Fatal(args...)
}

// Fatalf 有致命性错误，导致程序崩溃，记录日志，然后退出
func (zl *ZapLog) Fatalf(ctx context.Context, template string, args ...interface{}) {
	zl.logger.Sugar().Fatalf(template, args...)
}

// Error 错误日志
func (zl *ZapLog) Error(ctx context.Context, args ...interface{}) {
	zl.logger.Sugar().Error(args...)
}

// Errorf 错误日志
func (zl *ZapLog) Errorf(ctx context.Context, template string, args ...interface{}) {
	zl.logger.Sugar().Errorf(template, args...)
}

// Warn 警告日志
func (zl *ZapLog) Warn(ctx context.Context, args ...interface{}) {
	zl.logger.Sugar().Warn(args...)
}

// Warnf 警告日志
func (zl *ZapLog) Warnf(ctx context.Context, template string, args ...interface{}) {
	zl.logger.Sugar().Warnf(template, args...)
}

// Info 核心流程日志
func (zl *ZapLog) Info(ctx context.Context, args ...interface{}) {
	zl.logger.Sugar().Info(args...)
}

// Infof 核心流程日志
func (zl *ZapLog) Infof(ctx context.Context, template string, args ...interface{}) {
	zl.logger.Sugar().Infof(template, args...)
}

// Debug debug日志（调试日志）
func (zl *ZapLog) Debug(ctx context.Context, args ...interface{}) {
	zl.logger.Sugar().Debug(args...)
}

// Debugf debug日志（调试日志）
func (zl *ZapLog) Debugf(ctx context.Context, template string, args ...interface{}) {
	zl.logger.Sugar().Debugf(template, args...)
}

// Trace 粒度超细的，一般情况下我们使用不上
func (zl *ZapLog) Trace(ctx context.Context, args ...interface{}) {
	zl.logger.Sugar().Debug(args...)
}
