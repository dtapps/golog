package golog

import (
	"context"
	"fmt"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gourl"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"unicode/utf8"
)

// 记录日志
func (ag *ApiGorm) gormRecord(ctx context.Context, data GormApiLogModel) {
	if ag.gormConfig.stats == false {
		return
	}

	// OpenTelemetry链路追踪
	if ag.trace {
		tr := otel.Tracer("go.dtapp.net/golog", trace.WithInstrumentationVersion(Version))
		ctx, ag.span = tr.Start(ctx, "api")
		defer ag.span.End()
	}

	data.GoVersion = ag.config.GoVersion                         //【程序】GoVersion
	data.SdkVersion = ag.config.SdkVersion                       //【程序】SdkVersion
	data.SystemInfo = gojson.JsonEncodeNoError(ag.config.system) //【系统】SystemInfo

	if utf8.ValidString(data.ResponseBody) == false {
		data.ResponseBody = ""
	}

	// OpenTelemetry链路追踪
	if ag.trace {
		data.TraceID = ag.span.SpanContext().TraceID().String() // 跟踪编号
	}

	// 请求编号
	data.RequestID = gorequest.GetRequestIDContext(ctx)

	err := ag.gormClient.WithContext(ctx).
		Table(ag.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		slog.Error(fmt.Sprintf("记录接口日志错误：%s", err))
		slog.Error(fmt.Sprintf("记录接口日志数据：%+v", data))
	}
}

// 中间件
func (ag *ApiGorm) gormMiddleware(ctx context.Context, request gorequest.Response) {
	data := GormApiLogModel{
		RequestTime:        request.RequestTime,                              //【请求】时间
		RequestUri:         request.RequestUri,                               //【请求】链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,           //【请求】链接
		RequestApi:         gourl.UriParse(request.RequestUri).Path,          //【请求】接口
		RequestMethod:      request.RequestMethod,                            //【请求】方式
		RequestParams:      gojson.JsonEncodeNoError(request.RequestParams),  //【请求】参数
		RequestHeader:      gojson.JsonEncodeNoError(request.RequestHeader),  //【请求】头部
		ResponseHeader:     gojson.JsonEncodeNoError(request.ResponseHeader), //【返回】头部
		ResponseStatusCode: request.ResponseStatusCode,                       //【返回】状态码
		ResponseTime:       request.ResponseTime,                             //【返回】时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(request.ResponseBody))) //【返回】数据
		}
	}

	ag.gormRecord(ctx, data)
}

// 中间件
func (ag *ApiGorm) gormMiddlewareXml(ctx context.Context, request gorequest.Response) {
	data := GormApiLogModel{
		RequestTime:        request.RequestTime,                              //【请求】时间
		RequestUri:         request.RequestUri,                               //【请求】链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,           //【请求】链接
		RequestApi:         gourl.UriParse(request.RequestUri).Path,          //【请求】接口
		RequestMethod:      request.RequestMethod,                            //【请求】方式
		RequestParams:      gojson.JsonEncodeNoError(request.RequestParams),  //【请求】参数
		RequestHeader:      gojson.JsonEncodeNoError(request.RequestHeader),  //【请求】头部
		ResponseHeader:     gojson.JsonEncodeNoError(request.ResponseHeader), //【返回】头部
		ResponseStatusCode: request.ResponseStatusCode,                       //【返回】状态码
		ResponseTime:       request.ResponseTime,                             //【返回】时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.XmlEncodeNoError(gojson.XmlDecodeNoError(request.ResponseBody)) //【返回】内容
		}
	}

	ag.gormRecord(ctx, data)
}

// 中间件
func (ag *ApiGorm) gormMiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	data := GormApiLogModel{
		RequestTime:        request.RequestTime,                              //【请求】时间
		RequestUri:         request.RequestUri,                               //【请求】链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,           //【请求】链接
		RequestApi:         api,                                              //【请求】接口
		RequestMethod:      request.RequestMethod,                            //【请求】方式
		RequestParams:      gojson.JsonEncodeNoError(request.RequestParams),  //【请求】参数
		RequestHeader:      gojson.JsonEncodeNoError(request.RequestHeader),  //【请求】头部
		ResponseHeader:     gojson.JsonEncodeNoError(request.ResponseHeader), //【返回】头部
		ResponseStatusCode: request.ResponseStatusCode,                       //【返回】状态码
		ResponseTime:       request.ResponseTime,                             //【返回】时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(request.ResponseBody))) //【返回】数据
		}
	}

	ag.gormRecord(ctx, data)
}
