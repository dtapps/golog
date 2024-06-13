package golog

import (
	"context"
	"fmt"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gourl"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"log/slog"
	"unicode/utf8"
)

// 记录日志
func (ag *ApiGorm) gormRecord(ctx context.Context, data GormApiLogModel) {
	if ag.gormConfig.stats == false {
		return
	}

	if utf8.ValidString(data.ResponseBody) == false {
		data.ResponseBody = ""
	}

	// 跟踪编号
	data.TraceID = ag.TraceGetTraceID()

	// 请求编号
	data.RequestID = gorequest.GetRequestIDContext(ctx)

	// OpenTelemetry链路追踪
	ag.TraceSetAttributes(attribute.String("request.id", data.RequestID))
	ag.TraceSetAttributes(attribute.String("request.time", data.RequestTime.Format(gotime.DateTimeFormat)))
	ag.TraceSetAttributes(attribute.String("request.uri", data.RequestUri))
	ag.TraceSetAttributes(attribute.String("request.url", data.RequestUrl))
	ag.TraceSetAttributes(attribute.String("request.api", data.RequestApi))
	ag.TraceSetAttributes(attribute.String("request.method", data.RequestMethod))
	ag.TraceSetAttributes(attribute.String("request.params", data.RequestParams))
	ag.TraceSetAttributes(attribute.String("request.header", data.RequestHeader))
	ag.TraceSetAttributes(attribute.String("request.ip", data.RequestIP))
	ag.TraceSetAttributes(attribute.Int64("request.cost_time", data.RequestCostTime))
	ag.TraceSetAttributes(attribute.String("response.header", data.ResponseHeader))
	ag.TraceSetAttributes(attribute.Int("response.status_code", data.ResponseStatusCode))
	ag.TraceSetAttributes(attribute.String("response.body", data.ResponseBody))
	ag.TraceSetAttributes(attribute.String("response.time", data.ResponseTime.Format(gotime.DateTimeFormat)))

	err := ag.gormClient.WithContext(ctx).
		Table(ag.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		ag.TraceRecordError(err)
		ag.TraceSetStatus(codes.Error, err.Error())
		slog.Error(fmt.Sprintf("记录接口日志错误：%s", err))
		slog.Error(fmt.Sprintf("记录接口日志数据：%+v", data))
	}
}

// 中间件
func (ag *ApiGorm) gormMiddleware(ctx context.Context, request gorequest.Response) {
	data := GormApiLogModel{
		RequestTime:        request.RequestTime,                              // 请求时间
		RequestUri:         request.RequestUri,                               // 请求链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,           // 请求链接
		RequestApi:         gourl.UriParse(request.RequestUri).Path,          // 请求接口
		RequestMethod:      request.RequestMethod,                            // 请求方式
		RequestParams:      gojson.JsonEncodeNoError(request.RequestParams),  // 请求参数
		RequestHeader:      gojson.JsonEncodeNoError(request.RequestHeader),  // 请求头部
		RequestCostTime:    request.RequestCostTime,                          // 请求消耗时长
		ResponseHeader:     gojson.JsonEncodeNoError(request.ResponseHeader), // 响应头部
		ResponseStatusCode: request.ResponseStatusCode,                       // 响应状态码
		ResponseTime:       request.ResponseTime,                             // 响应时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(request.ResponseBody))) // 响应数据
		}
	}

	ag.gormRecord(ctx, data)
}

// 中间件
func (ag *ApiGorm) gormMiddlewareXml(ctx context.Context, request gorequest.Response) {
	data := GormApiLogModel{
		RequestTime:        request.RequestTime,                              // 请求时间
		RequestUri:         request.RequestUri,                               // 请求链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,           // 请求链接
		RequestApi:         gourl.UriParse(request.RequestUri).Path,          // 请求接口
		RequestMethod:      request.RequestMethod,                            // 请求方式
		RequestParams:      gojson.JsonEncodeNoError(request.RequestParams),  // 请求参数
		RequestHeader:      gojson.JsonEncodeNoError(request.RequestHeader),  // 请求头部
		RequestCostTime:    request.RequestCostTime,                          // 请求消耗时长
		ResponseHeader:     gojson.JsonEncodeNoError(request.ResponseHeader), // 响应头部
		ResponseStatusCode: request.ResponseStatusCode,                       // 响应状态码
		ResponseTime:       request.ResponseTime,                             // 响应时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.XmlEncodeNoError(gojson.XmlDecodeNoError(request.ResponseBody)) // 响应内容
		}
	}

	ag.gormRecord(ctx, data)
}

// 中间件
func (ag *ApiGorm) gormMiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	data := GormApiLogModel{
		RequestTime:        request.RequestTime,                              // 请求时间
		RequestUri:         request.RequestUri,                               // 请求链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,           // 请求链接
		RequestApi:         api,                                              // 请求接口
		RequestMethod:      request.RequestMethod,                            // 请求方式
		RequestParams:      gojson.JsonEncodeNoError(request.RequestParams),  // 请求参数
		RequestHeader:      gojson.JsonEncodeNoError(request.RequestHeader),  // 请求头部
		RequestCostTime:    request.RequestCostTime,                          // 请求消耗时长
		ResponseHeader:     gojson.JsonEncodeNoError(request.ResponseHeader), // 响应头部
		ResponseStatusCode: request.ResponseStatusCode,                       // 响应状态码
		ResponseTime:       request.ResponseTime,                             // 响应时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(request.ResponseBody))) // 响应数据
		}
	}

	ag.gormRecord(ctx, data)
}
