package golog

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// gormRecord 记录日志
func (gm *GinMongo) gormRecord(ctx context.Context, data ginMongoLog) {
	if gm.mongoConfig.stats == false {
		return
	}

	data.LogID = primitive.NewObjectID()           //【记录】编号
	data.SystemHostName = gm.config.systemHostname //【系统】主机名
	data.SystemInsideIP = gm.config.systemInsideIP //【系统】内网IP
	data.GoVersion = gm.config.goVersion           //【程序】Go版本
	data.SdkVersion = gm.config.sdkVersion         //【程序】Sdk版本
	data.SystemVersion = gm.config.systemVersion   //【程序】System版本
	data.SystemOs = gm.config.systemOs             //【系统】类型
	data.SystemArch = gm.config.systemKernel       //【系统】架构
	data.SystemUpTime = gm.config.systemUpTime     //【系统】运行时间
	data.SystemBootTime = gm.config.systemBootTime //【系统】开机时间
	data.CpuCores = gm.config.cpuCores             //【CPU】核数
	data.CpuModelName = gm.config.cpuModelName     //【CPU】型号名称
	data.CpuMhz = gm.config.cpuMhz                 //【CPU】兆赫

	_, err := gm.mongoClient.Database(gm.mongoConfig.databaseName).
		Collection(gm.mongoConfig.collectionName).
		InsertOne(ctx, data)
	if err != nil {
		if gm.slog.status {
			gm.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录接口日志错误：%s", err))
		}
		if gm.slog.status {
			gm.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录接口日志数据：%+v", data))
		}
	}
}

func (gm *GinMongo) recordJson(ginCtx *gin.Context, requestTime time.Time, requestBody gorequest.Params, responseTime time.Time, responseCode int, responseBody any, costTime int64, requestIp string) {

	data := ginMongoLog{
		TraceID:       gotrace_id.GetGinTraceId(ginCtx),                             //【系统】跟踪编号
		LogTime:       primitive.NewDateTimeFromTime(requestTime),                   //【记录】时间
		RequestTime:   gotime.SetCurrent(requestTime).Format(),                      //【请求】时间
		RequestURL:    ginCtx.Request.RequestURI,                                    //【请求】链接
		RequestApi:    gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】接口
		RequestMethod: ginCtx.Request.Method,                                        //【请求】方式
		RequestProto:  ginCtx.Request.Proto,                                         //【请求】协议
		RequestBody:   requestBody,                                                  //【请求】参数
		RequestIP:     requestIp,                                                    //【请求】客户端IP
		RequestHeader: ginCtx.Request.Header,                                        //【请求】头部
		ResponseTime:  gotime.SetCurrent(responseTime).Format(),                     //【返回】时间
		ResponseCode:  responseCode,                                                 //【返回】状态码
		ResponseData:  responseBody,                                                 //【返回】数据
		CostTime:      costTime,                                                     //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】链接
	}

	gm.gormRecord(ginCtx, data)
}
