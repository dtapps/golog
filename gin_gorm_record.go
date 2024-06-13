package golog

import (
	"context"
	"fmt"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gotime"
	"go.opentelemetry.io/otel/attribute"
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

	// OpenTelemetry链路追踪
	gg.TraceSetAttributes(attribute.String("request.id", data.RequestID))
	gg.TraceSetAttributes(attribute.String("request.time", data.RequestTime.Format(gotime.DateTimeFormat)))
	gg.TraceSetAttributes(attribute.String("request.host", data.RequestHost))
	gg.TraceSetAttributes(attribute.String("request.path", data.RequestPath))
	gg.TraceSetAttributes(attribute.String("request.query", data.RequestQuery))
	gg.TraceSetAttributes(attribute.String("request.method", data.RequestMethod))
	gg.TraceSetAttributes(attribute.String("request.scheme", data.RequestScheme))
	gg.TraceSetAttributes(attribute.String("request.content_type", data.RequestContentType))
	gg.TraceSetAttributes(attribute.String("request.body", data.RequestBody))
	gg.TraceSetAttributes(attribute.String("request.client_ip", data.RequestClientIP))
	gg.TraceSetAttributes(attribute.String("request.user_agent", data.RequestClientIP))
	gg.TraceSetAttributes(attribute.String("request.header", data.RequestHeader))
	gg.TraceSetAttributes(attribute.Int64("request.cost_time", data.RequestCostTime))
	gg.TraceSetAttributes(attribute.String("response.time", data.ResponseTime.Format(gotime.DateTimeFormat)))
	gg.TraceSetAttributes(attribute.String("response.header", data.ResponseHeader))
	gg.TraceSetAttributes(attribute.Int("response.status_code", data.ResponseStatusCode))
	gg.TraceSetAttributes(attribute.String("response.body", data.ResponseBody))
	gg.TraceSetAttributes(attribute.String("system_info", data.SystemInfo))

	err := gg.gormClient.WithContext(ctx).
		Table(gg.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		gg.TraceRecordError(err)
		gg.TraceSetStatus(codes.Error, err.Error())
		slog.Error(fmt.Sprintf("记录接口日志错误：%s", err))
		slog.Error(fmt.Sprintf("记录接口日志数据：%+v", data))
	}
}
