package golog

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"time"
)

// gormRecord 记录日志
func (gg *GinGorm) gormRecord(ctx context.Context, data ginGormLog) {
	if gg.gormConfig.stats == false {
		return
	}

	data.SystemHostName = gg.config.systemHostname //【系统】主机名
	data.SystemInsideIP = gg.config.systemInsideIP //【系统】内网ip
	data.GoVersion = gg.config.goVersion           //【程序】Go版本
	data.SdkVersion = gg.config.sdkVersion         //【程序】Sdk版本
	data.SystemVersion = gg.config.systemVersion   //【程序】System版本
	data.SystemOs = gg.config.systemOs             //【系统】类型
	data.SystemArch = gg.config.systemKernel       //【系统】架构
	data.SystemUpTime = gg.config.systemUpTime     //【系统】运行时间
	data.SystemBootTime = gg.config.systemBootTime //【系统】开机时间
	data.CpuCores = gg.config.cpuCores             //【CPU】核数
	data.CpuModelName = gg.config.cpuModelName     //【CPU】型号名称
	data.CpuMhz = gg.config.cpuMhz                 //【CPU】兆赫

	err := gg.gormClient.WithContext(ctx).Table(gg.gormConfig.tableName).Create(&data).Error
	if err != nil {
		if gg.slog.status {
			gg.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录接口日志错误：%s", err))
		}
		if gg.slog.status {
			gg.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录接口日志数据：%+v", data))
		}
	}
}

func (gg *GinGorm) recordJson(ginCtx *gin.Context, requestTime time.Time, requestBody gorequest.Params, responseTime time.Time, responseCode int, responseBody string, costTime int64, requestIp string) {

	data := ginGormLog{
		TraceID:       gotrace_id.GetGinTraceId(ginCtx),                             //【系统】跟踪编号
		RequestTime:   requestTime,                                                  //【请求】时间
		RequestURL:    ginCtx.Request.RequestURI,                                    //【请求】链接
		RequestApi:    gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】接口
		RequestMethod: ginCtx.Request.Method,                                        //【请求】方式
		RequestProto:  ginCtx.Request.Proto,                                         //【请求】协议
		RequestBody:   gojson.JsonEncodeNoError(requestBody),                        //【请求】参数
		RequestIP:     requestIp,                                                    //【请求】客户端IP
		RequestHeader: gojson.JsonEncodeNoError(ginCtx.Request.Header),              //【请求】头部
		ResponseTime:  responseTime,                                                 //【返回】时间
		ResponseCode:  responseCode,                                                 //【返回】状态码
		ResponseData:  responseBody,                                                 //【返回】数据
		CostTime:      costTime,                                                     //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】链接
	}

	gg.gormRecord(ginCtx, data)
}
