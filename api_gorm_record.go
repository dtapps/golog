package golog

import (
	"context"
	"fmt"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"unicode/utf8"
)

// 记录日志
func (ag *ApiGorm) gormRecord(ctx context.Context, data apiGormLog) {
	if ag.gormConfig.stats == false {
		return
	}

	if utf8.ValidString(data.ResponseBody) == false {
		data.ResponseBody = ""
	}

	data.TraceID = gotrace_id.GetTraceIdContext(ctx) //【记录】跟踪编号
	data.SystemHostName = ag.config.systemHostname   //【系统】主机名
	data.SystemInsideIP = ag.config.systemInsideIP   //【系统】内网IP
	data.GoVersion = ag.config.goVersion             //【程序】Go版本
	data.SdkVersion = ag.config.sdkVersion           //【程序】Sdk版本
	data.SystemVersion = ag.config.systemVersion     //【程序】System版本
	data.RequestIP = ag.config.systemOutsideIP       //【请求】请求Ip
	data.SystemOs = ag.config.systemOs               //【系统】类型
	data.SystemArch = ag.config.systemKernel         //【系统】架构
	data.SystemUpTime = ag.config.systemUpTime       //【系统】运行时间
	data.SystemBootTime = ag.config.systemBootTime   //【系统】开机时间
	data.CpuCores = ag.config.cpuCores               //【CPU】核数
	data.CpuModelName = ag.config.cpuModelName       //【CPU】型号名称
	data.CpuMhz = ag.config.cpuMhz                   //【CPU】兆赫

	err := ag.gormClient.WithContext(ctx).Table(ag.gormConfig.tableName).Create(&data).Error
	if err != nil {
		if ag.slog.status {
			ag.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录接口日志错误：%s", err))
		}
		if ag.slog.status {
			ag.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录接口日志数据：%+v", data))
		}
	}
}

// 中间件
func (ag *ApiGorm) gormMiddleware(ctx context.Context, request gorequest.Response) {
	data := apiGormLog{
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
	data := apiGormLog{
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
	data := apiGormLog{
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
