package golog

import (
	"context"
	"fmt"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 记录日志
func (am *ApiMongo) mongoRecord(ctx context.Context, data apiMongolLog) {
	if am.mongoConfig.stats == false {
		return
	}

	data.LogID = primitive.NewObjectID()             //【记录】编号
	data.TraceID = gotrace_id.GetTraceIdContext(ctx) //【记录】跟踪编号
	data.SystemHostName = am.config.systemHostname   //【系统】主机名
	data.SystemInsideIP = am.config.systemInsideIP   //【系统】内网IP
	data.GoVersion = am.config.goVersion             //【系统】Go版本
	data.SdkVersion = am.config.sdkVersion           //【系统】Sdk版本
	data.SystemVersion = am.config.systemVersion     //【系统】System版本
	data.RequestIP = am.config.systemOutsideIP       //【请求】请求Ip
	data.SystemOs = am.config.systemOs               //【系统】类型
	data.SystemArch = am.config.systemKernel         //【系统】架构
	data.SystemUpTime = am.config.systemUpTime       //【系统】运行时间
	data.SystemBootTime = am.config.systemBootTime   //【系统】开机时间
	data.CpuCores = am.config.cpuCores               //【CPU】核数
	data.CpuModelName = am.config.cpuModelName       //【CPU】型号名称
	data.CpuMhz = am.config.cpuMhz                   //【CPU】兆赫

	_, err := am.mongoClient.Database(am.mongoConfig.databaseName).
		Collection(am.mongoConfig.collectionName).
		InsertOne(ctx, data)
	if err != nil {
		if am.slog.status {
			am.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录接口日志错误：%s", err))
		}
		if am.slog.status {
			am.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录接口日志数据：%+v", data))
		}
	}
}

// 中间件
func (am *ApiMongo) mongoMiddleware(ctx context.Context, request gorequest.Response) {
	data := apiMongolLog{
		LogTime:            primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:        gotime.SetCurrent(request.RequestTime).Format(),    //【请求】时间
		RequestUri:         request.RequestUri,                                 //【请求】链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:         gourl.UriParse(request.RequestUri).Path,            //【请求】接口
		RequestMethod:      request.RequestMethod,                              //【请求】方式
		RequestParams:      request.RequestParams,                              //【请求】参数
		RequestHeader:      request.RequestHeader,                              //【请求】头部
		ResponseHeader:     request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode: request.ResponseStatusCode,                         //【返回】状态码
		ResponseTime:       gotime.SetCurrent(request.ResponseTime).Format(),   //【返回】时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.JsonDecodeNoError(string(request.ResponseBody)) //【返回】内容
		}
	}

	am.mongoRecord(ctx, data)
}

// 中间件
func (am *ApiMongo) mongoMiddlewareXml(ctx context.Context, request gorequest.Response) {
	data := apiMongolLog{
		LogTime:            primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:        gotime.SetCurrent(request.RequestTime).Format(),    //【请求】时间
		RequestUri:         request.RequestUri,                                 //【请求】链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:         gourl.UriParse(request.RequestUri).Path,            //【请求】接口
		RequestMethod:      request.RequestMethod,                              //【请求】方式
		RequestParams:      request.RequestParams,                              //【请求】参数
		RequestHeader:      request.RequestHeader,                              //【请求】头部
		ResponseHeader:     request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode: request.ResponseStatusCode,                         //【返回】状态码
		ResponseTime:       gotime.SetCurrent(request.ResponseTime).Format(),   //【返回】时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.XmlDecodeNoError(request.ResponseBody) //【返回】内容
		}
	}

	am.mongoRecord(ctx, data)
}

// 中间件
func (am *ApiMongo) mongoMiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	data := apiMongolLog{
		LogTime:            primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:        gotime.SetCurrent(request.RequestTime).Format(),    //【请求】时间
		RequestUri:         request.RequestUri,                                 //【请求】链接
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:         api,                                                //【请求】接口
		RequestMethod:      request.RequestMethod,                              //【请求】方式
		RequestParams:      request.RequestParams,                              //【请求】参数
		RequestHeader:      request.RequestHeader,                              //【请求】头部
		ResponseHeader:     request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode: request.ResponseStatusCode,                         //【返回】状态码
		ResponseTime:       gotime.SetCurrent(request.ResponseTime).Format(),   //【返回】时间
	}
	if !request.HeaderIsImg() {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = gojson.JsonDecodeNoError(string(request.ResponseBody)) //【返回】内容
		}
	}

	am.mongoRecord(ctx, data)
}
