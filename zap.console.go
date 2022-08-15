package golog

import (
	"context"
	"go.dtapp.net/gotrace_id"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type ZapConsole struct {
	logger *zap.Logger // 日志服务
}

func NewZapConsole() *ZapConsole {

	zc := &ZapConsole{}

	zc.logger = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{}),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	))

	return zc
}

// 跟踪编号
func (zc *ZapConsole) getTraceId(ctx context.Context) zap.Field {
	traceId := gotrace_id.GetTraceIdContext(ctx)
	if traceId == "" {
		return zap.Field{}
	} else {
		return zap.String("trace_id", traceId)
	}
}

// Print 打印
func (zc *ZapConsole) Print(ctx context.Context, args ...interface{}) {
	zc.logger.Sugar().Info(args, zc.getTraceId(ctx))
}

// Printf 打印
func (zc *ZapConsole) Printf(ctx context.Context, template string, args ...interface{}) {
	zc.logger.Sugar().Infof(template, args, zc.getTraceId(ctx))
}
