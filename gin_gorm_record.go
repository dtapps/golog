package golog

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/dorm"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gourl"
	"log"
	"time"
)

// 模型
type ginPostgresqlLog struct {
	LogID           uint      `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`       //【记录】编号
	TraceID         string    `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`        //【系统】跟踪编号
	RequestTime     time.Time `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`      //【请求】时间
	RequestUri      string    `gorm:"comment:【请求】请求链接 域名+路径+参数" json:"request_uri,omitempty"`  //【请求】请求链接 域名+路径+参数
	RequestURL      string    `gorm:"comment:【请求】请求链接 域名+路径" json:"request_url,omitempty"`     //【请求】请求链接 域名+路径
	RequestApi      string    `gorm:"index;comment:【请求】请求接口 路径" json:"request_api,omitempty"`  //【请求】请求接口 路径
	RequestMethod   string    `gorm:"index;comment:【请求】请求方式" json:"request_method,omitempty"`  //【请求】请求方式
	RequestProto    string    `gorm:"comment:【请求】请求协议" json:"request_proto,omitempty"`         //【请求】请求协议
	RequestUa       string    `gorm:"comment:【请求】请求UA" json:"request_ua,omitempty"`            //【请求】请求UA
	RequestReferer  string    `gorm:"comment:【请求】请求referer" json:"request_referer,omitempty"`  //【请求】请求referer
	RequestBody     string    `gorm:"comment:【请求】请求主体" json:"request_body,omitempty"`          //【请求】请求主体
	RequestUrlQuery string    `gorm:"comment:【请求】请求URL参数" json:"request_url_query,omitempty"`  //【请求】请求URL参数
	RequestIP       string    `gorm:"index;comment:【请求】请求客户端IP" json:"request_ip,omitempty"`   //【请求】请求客户端IP
	RequestHeader   string    `gorm:"comment:【请求】请求头" json:"request_header,omitempty"`         //【请求】请求头
	ResponseTime    time.Time `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`     //【返回】时间
	ResponseCode    int       `gorm:"index;comment:【返回】状态码" json:"response_code,omitempty"`    //【返回】状态码
	ResponseMsg     string    `gorm:"comment:【返回】描述" json:"response_msg,omitempty"`            //【返回】描述
	ResponseData    string    `gorm:"comment:【返回】数据" json:"response_data,omitempty"`           //【返回】数据
	CostTime        int64     `gorm:"comment:【系统】花费时间" json:"cost_time,omitempty"`             //【系统】花费时间
	SystemHostName  string    `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"` //【系统】主机名
	SystemInsideIP  string    `gorm:"comment:【系统】内网ip" json:"system_inside_ip,omitempty"`      //【系统】内网ip
	SystemOs        string    `gorm:"index;comment:【系统】系统类型" json:"system_os,omitempty"`       //【系统】系统类型
	SystemArch      string    `gorm:"index;comment:【系统】系统架构" json:"system_arch,omitempty"`     //【系统】系统架构
	GoVersion       string    `gorm:"comment:【程序】Go版本" json:"go_version,omitempty"`            //【程序】Go版本
	SdkVersion      string    `gorm:"comment:【程序】Sdk版本" json:"sdk_version,omitempty"`          //【程序】Sdk版本
}

// 创建模型
func (gg *GinGorm) gormAutoMigrate(ctx context.Context) {
	err := gg.gormClient.GetDb().Table(gg.gormConfig.tableName).AutoMigrate(&ginPostgresqlLog{})
	if err != nil {
		log.Printf("创建模型：%s\n", err)
	}
}

// gormRecord 记录日志
func (gg *GinGorm) gormRecord(data ginPostgresqlLog) {

	data.SystemHostName = gg.config.systemHostname //【系统】主机名
	data.SystemInsideIP = gg.config.systemInsideIp //【系统】内网ip
	data.GoVersion = gg.config.goVersion           //【程序】Go版本
	data.SdkVersion = gg.config.sdkVersion         //【程序】Sdk版本
	data.SystemOs = gg.config.systemOs             //【系统】系统类型
	data.SystemArch = gg.config.systemKernel       //【系统】系统架构

	err := gg.gormClient.GetDb().Table(gg.gormConfig.tableName).Create(&data).Error
	if err != nil {
		log.Printf("记录框架日志错误：%s\n", err)
		log.Printf("记录框架日志数据：%+v\n", data)
	}
}

func (gg *GinGorm) recordJson(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody gorequest.Params, responseCode int, responseBody string, startTime, endTime int64, requestIp string) {

	data := ginPostgresqlLog{
		TraceID:         traceId,                                                      //【系统】跟踪编号
		RequestTime:     requestTime,                                                  //【请求】时间
		RequestURL:      ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:      gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:   ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:    ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:       ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:  ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestBody:     dorm.JsonEncodeNoError(requestBody),                          //【请求】请求主体
		RequestUrlQuery: dorm.JsonEncodeNoError(ginCtx.Request.URL.Query()),           //【请求】请求URL参数
		RequestIP:       requestIp,                                                    //【请求】请求客户端IP
		RequestHeader:   dorm.JsonEncodeNoError(ginCtx.Request.Header),                //【请求】请求头
		ResponseTime:    gotime.Current().Time,                                        //【返回】时间
		ResponseCode:    responseCode,                                                 //【返回】状态码
		ResponseData:    responseBody,                                                 //【返回】数据
		CostTime:        endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	gg.gormRecord(data)
}
