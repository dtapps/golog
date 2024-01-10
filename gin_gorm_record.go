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

// 模型
type ginGormLog struct {
	LogID          int64     `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`                   //【记录】编号
	TraceID        string    `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`                    //【系统】跟踪编号
	RequestTime    time.Time `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`                  //【请求】时间
	RequestUri     string    `gorm:"comment:【请求】链接 域名+路径+参数" json:"request_uri,omitempty"`                //【请求】链接 域名+路径+参数
	RequestURL     string    `gorm:"comment:【请求】链接 域名+路径" json:"request_url,omitempty"`                   //【请求】链接 域名+路径
	RequestApi     string    `gorm:"index;comment:【请求】接口" json:"request_api,omitempty"`                   //【请求】接口
	RequestMethod  string    `gorm:"index;comment:【请求】方式" json:"request_method,omitempty"`                //【请求】方式
	RequestProto   string    `gorm:"comment:【请求】协议" json:"request_proto,omitempty"`                       //【请求】协议
	RequestBody    string    `gorm:"comment:【请求】参数" json:"request_body,omitempty"`                        //【请求】参数
	RequestIP      string    `gorm:"default:0.0.0.0;index;comment:【请求】客户端IP" json:"request_ip,omitempty"` //【请求】客户端IP
	RequestHeader  string    `gorm:"comment:【请求】头部" json:"request_header,omitempty"`                      //【请求】头部
	ResponseTime   time.Time `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`                 //【返回】时间
	ResponseCode   int       `gorm:"comment:【返回】状态码" json:"response_code,omitempty"`                      //【返回】状态码
	ResponseData   string    `gorm:"comment:【返回】数据" json:"response_data,omitempty"`                       //【返回】数据
	CostTime       int64     `gorm:"comment:【系统】花费时间" json:"cost_time,omitempty"`                         //【系统】花费时间
	SystemHostName string    `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"`             //【系统】主机名
	SystemInsideIP string    `gorm:"default:0.0.0.0;comment:【系统】内网IP" json:"system_inside_ip,omitempty"`  //【系统】内网IP
	SystemOs       string    `gorm:"comment:【系统】类型" json:"system_os,omitempty"`                           //【系统】类型
	SystemArch     string    `gorm:"comment:【系统】架构" json:"system_arch,omitempty"`                         //【系统】架构
	SystemUpTime   uint64    `gorm:"comment:【系统】运行时间" json:"system_up_time,omitempty"`                    //【系统】运行时间
	SystemBootTime uint64    `gorm:"comment:【系统】开机时间" json:"system_boot_time,omitempty"`                  //【系统】开机时间
	GoVersion      string    `gorm:"comment:【程序】Go版本" json:"go_version,omitempty"`                        //【程序】Go版本
	SdkVersion     string    `gorm:"comment:【程序】Sdk版本" json:"sdk_version,omitempty"`                      //【程序】Sdk版本
	SystemVersion  string    `gorm:"comment:【程序】System版本" json:"system_version,omitempty"`                //【程序】System版本
	CpuCores       int       `gorm:"comment:【CPU】核数" json:"cpu_cores,omitempty"`                          //【CPU】核数
	CpuModelName   string    `gorm:"comment:【CPU】型号名称" json:"cpu_model_name,omitempty"`                   //【CPU】型号名称
	CpuMhz         float64   `gorm:"comment:【CPU】兆赫" json:"cpu_mhz,omitempty"`                            //【CPU】兆赫
}

// 创建模型
func (gg *GinGorm) gormAutoMigrate(ctx context.Context) {
	if gg.gormConfig.stats == false {
		return
	}

	err := gg.gormClient.WithContext(ctx).Table(gg.gormConfig.tableName).AutoMigrate(&ginGormLog{})
	if err != nil {
		if gg.slog.status {
			gg.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("创建模型：%s", err))
		}
	}
}

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
