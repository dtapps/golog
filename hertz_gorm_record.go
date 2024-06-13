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
func (hg *HertzGorm) gormRecord(ctx context.Context, data GormHertzLogModel) {
	if hg.gormConfig.stats == false {
		return
	}
	data.GoVersion = hg.config.GoVersion                         //【程序】GoVersion
	data.SdkVersion = hg.config.SdkVersion                       //【程序】SdkVersion
	data.SystemInfo = gojson.JsonEncodeNoError(hg.config.system) //【系统】SystemInfo

	// OpenTelemetry链路追踪
	hg.TraceSetAttributes(attribute.String("request.id", data.RequestID))
	hg.TraceSetAttributes(attribute.String("request.time", data.RequestTime.Format(gotime.DateTimeFormat)))
	hg.TraceSetAttributes(attribute.String("request.host", data.RequestHost))
	hg.TraceSetAttributes(attribute.String("request.path", data.RequestPath))
	hg.TraceSetAttributes(attribute.String("request.query", data.RequestQuery))
	hg.TraceSetAttributes(attribute.String("request.method", data.RequestMethod))
	hg.TraceSetAttributes(attribute.String("request.scheme", data.RequestScheme))
	hg.TraceSetAttributes(attribute.String("request.content_type", data.RequestContentType))
	hg.TraceSetAttributes(attribute.String("request.body", data.RequestBody))
	hg.TraceSetAttributes(attribute.String("request.client_ip", data.RequestClientIP))
	hg.TraceSetAttributes(attribute.String("request.user_agent", data.RequestClientIP))
	hg.TraceSetAttributes(attribute.String("request.header", data.RequestHeader))
	hg.TraceSetAttributes(attribute.Int64("request.cost_time", data.RequestCostTime))
	hg.TraceSetAttributes(attribute.String("response.time", data.ResponseTime.Format(gotime.DateTimeFormat)))
	hg.TraceSetAttributes(attribute.String("response.header", data.ResponseHeader))
	hg.TraceSetAttributes(attribute.Int("response.status_code", data.ResponseStatusCode))
	hg.TraceSetAttributes(attribute.String("response.body", data.ResponseBody))
	hg.TraceSetAttributes(attribute.String("system_info", data.SystemInfo))

	err := hg.gormClient.WithContext(ctx).
		Table(hg.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		hg.TraceSetStatus(codes.Error, err.Error())
		hg.TraceRecordError(err)
		slog.Error(fmt.Sprintf("记录接口日志错误：%s", err))
		slog.Error(fmt.Sprintf("记录接口日志数据：%+v", data))
	}
}
