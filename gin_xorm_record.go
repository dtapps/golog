package golog

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"log"
	"time"
)

// 模型
type ginXormLog struct {
	LogID          int64     `xorm:"pk;comment('【记录】编号')" json:"log_id,omitempty"`               //【记录】编号
	TraceID        string    `xorm:"index;comment('【系统】跟踪编号')" json:"trace_id,omitempty"`        //【系统】跟踪编号
	RequestTime    time.Time `xorm:"index;comment('【请求】时间')" json:"request_time,omitempty"`      //【请求】时间
	RequestUri     string    `xorm:"comment('【请求】链接 域名+路径+参数')" json:"request_uri,omitempty"`    //【请求】链接 域名+路径+参数
	RequestURL     string    `xorm:"comment('【请求】链接 域名+路径')" json:"request_url,omitempty"`       //【请求】链接 域名+路径
	RequestApi     string    `xorm:"index;comment('【请求】接口')" json:"request_api,omitempty"`       //【请求】接口
	RequestMethod  string    `xorm:"index;comment('【请求】方式')" json:"request_method,omitempty"`    //【请求】方式
	RequestProto   string    `xorm:"comment('【请求】协议')" json:"request_proto,omitempty"`           //【请求】协议
	RequestBody    string    `xorm:"comment('【请求】参数')" json:"request_body,omitempty"`            //【请求】参数
	RequestIP      string    `xorm:"index;comment('【请求】客户端IP')" json:"request_ip,omitempty"`     //【请求】客户端IP
	RequestHeader  string    `xorm:"comment('【请求】头部')" json:"request_header,omitempty"`          //【请求】头部
	ResponseTime   time.Time `xorm:"index;comment('【返回】时间')" json:"response_time,omitempty"`     //【返回】时间
	ResponseCode   int       `xorm:"index;comment('【返回】状态码')" json:"response_code,omitempty"`    //【返回】状态码
	ResponseData   string    `xorm:"comment('【返回】数据')" json:"response_data,omitempty"`           //【返回】数据
	CostTime       int64     `xorm:"comment('【系统】花费时间')" json:"cost_time,omitempty"`             //【系统】花费时间
	SystemHostName string    `xorm:"index;comment('【系统】主机名')" json:"system_host_name,omitempty"` //【系统】主机名
	SystemInsideIP string    `xorm:"comment('【系统】内网IP')" json:"system_inside_ip,omitempty"`      //【系统】内网IP
	SystemOs       string    `xorm:"index;comment('【系统】系统类型')" json:"system_os,omitempty"`       //【系统】系统类型
	SystemArch     string    `xorm:"index;comment('【记系统录】系统架构')" json:"system_arch,omitempty"`   //【系统】系统架构
	GoVersion      string    `xorm:"index;comment('【程序】Go版本')" json:"go_version,omitempty"`      //【程序】Go版本
	SdkVersion     string    `xorm:"index;comment('【程序】Sdk版本')" json:"sdk_version,omitempty"`    //【程序】Sdk版本
}

// 创建模型
func (gg *GinXorm) xormSync(ctx context.Context) {
	if gg.xormConfig.stats == false {
		return
	}

	err := gg.xormClient.Table(gg.xormConfig.tableName).Sync(&ginXormLog{})
	if err != nil {
		log.Printf("创建模型：%s\n", err)
	}
}

// xormRecord 记录日志
func (gg *GinXorm) xormRecord(data ginXormLog) {
	if gg.xormConfig.stats == false {
		return
	}

	data.SystemHostName = gg.config.systemHostname //【系统】主机名
	data.SystemInsideIP = gg.config.systemInsideIp //【系统】内网ip
	data.GoVersion = gg.config.goVersion           //【程序】Go版本
	data.SdkVersion = gg.config.sdkVersion         //【程序】Sdk版本
	data.SystemOs = gg.config.systemOs             //【系统】系统类型
	data.SystemArch = gg.config.systemKernel       //【系统】系统架构

	_, err := gg.xormClient.Table(gg.xormConfig.tableName).Insert(&data)
	if err != nil {
		log.Printf("记录框架日志错误：%s\n", err)
		log.Printf("记录框架日志数据：%+v\n", data)
	}
}

func (gg *GinXorm) recordJson(ginCtx *gin.Context, requestTime time.Time, requestBody gorequest.Params, responseTime time.Time, responseCode int, responseBody string, costTime int64, requestIp string) {

	data := ginXormLog{
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

	gg.xormRecord(data)
}
