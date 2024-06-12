package golog

import (
	"context"
	"fmt"
	"go.dtapp.net/gojson"
	"go.opentelemetry.io/otel/codes"
	"log/slog"
)

// gormRecord 记录日志
func (gg *GinGorm) gormRecord(ctx context.Context, data GormGinLogModel) {
	if gg.gormConfig.stats == false {
		return
	}
	data.GoVersion = gg.config.GoVersion                         //【程序】GoVersion
	data.SdkVersion = gg.config.SdkVersion                       //【程序】SdkVersion
	data.SystemInfo = gojson.JsonEncodeNoError(gg.config.system) //【系统】SystemInfo

	err := gg.gormClient.WithContext(ctx).
		Table(gg.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		gg.TraceSetStatus(codes.Error, err.Error())
		slog.Error(fmt.Sprintf("记录接口日志错误：%s", err))
		slog.Error(fmt.Sprintf("记录接口日志数据：%+v", data))
	}
}
