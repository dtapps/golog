package golog

import (
	"context"
	"fmt"
	"go.dtapp.net/gojson"
	"go.opentelemetry.io/otel/codes"
	"log/slog"
)

// gormRecord 记录日志
func (hg *HertzGorm) gormRecord(ctx context.Context, data GormHertzLogModel) {
	if hg.gormConfig.stats == false {
		return
	}
	data.GoVersion = hg.config.GoVersion                         //【程序】GoVersion
	data.SdkVersion = hg.config.SdkVersion                       //【程序】SdkVersion
	data.SystemInfo = gojson.JsonEncodeNoError(hg.config.system) //【系统】SystemInfo

	err := hg.gormClient.WithContext(ctx).
		Table(hg.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		hg.TraceSetStatus(codes.Error, err.Error())
		slog.Error(fmt.Sprintf("记录接口日志错误：%s", err))
		slog.Error(fmt.Sprintf("记录接口日志数据：%+v", data))
	}
}
