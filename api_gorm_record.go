package golog

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"log"
	"time"
	"unicode/utf8"
)

// 模型
type apiGormLog struct {
	LogID              int64     `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`       //【记录】编号
	TraceID            string    `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`        //【系统】跟踪编号
	RequestTime        time.Time `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`      //【请求】时间
	RequestUri         string    `gorm:"comment:【请求】链接" json:"request_uri,omitempty"`             //【请求】链接
	RequestUrl         string    `gorm:"comment:【请求】链接" json:"request_url,omitempty"`             //【请求】链接
	RequestApi         string    `gorm:"index;comment:【请求】接口" json:"request_api,omitempty"`       //【请求】接口
	RequestMethod      string    `gorm:"index;comment:【请求】方式" json:"request_method,omitempty"`    //【请求】方式
	RequestParams      string    `gorm:"comment:【请求】参数" json:"request_params,omitempty"`          //【请求】参数
	RequestHeader      string    `gorm:"comment:【请求】头部" json:"request_header,omitempty"`          //【请求】头部
	RequestIp          string    `gorm:"comment:【请求】请求IP" json:"request_ip,omitempty"`            //【请求】请求IP
	ResponseHeader     string    `gorm:"comment:【返回】头部" json:"response_header,omitempty"`         //【返回】头部
	ResponseStatusCode int       `gorm:"comment:【返回】状态码" json:"response_status_code,omitempty"`   //【返回】状态码
	ResponseBody       string    `gorm:"comment:【返回】数据" json:"response_body,omitempty"`           //【返回】数据
	ResponseTime       time.Time `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`     //【返回】时间
	SystemHostName     string    `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"` //【系统】主机名
	SystemInsideIp     string    `gorm:"comment:【系统】内网IP" json:"system_inside_ip,omitempty"`      //【系统】内网IP
	SystemOs           string    `gorm:"comment:【系统】类型" json:"system_os,omitempty"`               //【系统】类型
	SystemArch         string    `gorm:"comment:【系统】架构" json:"system_arch,omitempty"`             //【系统】架构
	SystemUpTime       uint64    `gorm:"comment:【系统】运行时间" json:"system_up_time,omitempty"`        //【系统】运行时间
	SystemBootTime     uint64    `gorm:"comment:【系统】开机时间" json:"system_boot_time,omitempty"`      //【系统】开机时间
	GoVersion          string    `gorm:"comment:【程序】Go版本" json:"go_version,omitempty"`            //【程序】Go版本
	SdkVersion         string    `gorm:"comment:【程序】Sdk版本" json:"sdk_version,omitempty"`          //【程序】Sdk版本
	SystemVersion      string    `gorm:"comment:【程序】System版本" json:"system_version,omitempty"`    //【程序】System版本
	CpuCores           int       `gorm:"comment:【CPU】核数" json:"cpu_cores,omitempty"`              //【CPU】核数
	CpuModelName       string    `gorm:"comment:【CPU】型号名称" json:"cpu_model_name,omitempty"`       //【CPU】型号名称
	CpuMhz             float64   `gorm:"comment:【CPU】兆赫" json:"cpu_mhz,omitempty"`                //【CPU】兆赫
}

// 创建模型
func (ag *ApiGorm) gormAutoMigrate(ctx context.Context) {
	if ag.gormConfig.stats == false {
		return
	}

	err := ag.gormClient.Table(ag.gormConfig.tableName).AutoMigrate(&apiGormLog{})
	if err != nil {
		log.Printf("创建模型：%s", err)
	}
}

// 记录日志
func (ag *ApiGorm) gormRecord(ctx context.Context, data apiGormLog) {
	if ag.gormConfig.stats == false {
		return
	}

	if utf8.ValidString(data.ResponseBody) == false {
		data.ResponseBody = ""
	}

	data.SystemHostName = ag.config.systemHostname   //【系统】主机名
	data.SystemInsideIp = ag.config.systemInsideIp   //【系统】内网ip
	data.GoVersion = ag.config.goVersion             //【程序】Go版本
	data.SdkVersion = ag.config.sdkVersion           //【程序】Sdk版本
	data.SystemVersion = ag.config.systemVersion     //【程序】System版本
	data.TraceID = gotrace_id.GetTraceIdContext(ctx) //【记录】跟踪编号
	data.RequestIp = ag.config.systemOutsideIp       //【请求】请求Ip
	data.SystemOs = ag.config.systemOs               //【系统】类型
	data.SystemArch = ag.config.systemKernel         //【系统】架构
	data.SystemUpTime = ag.config.systemUpTime       //【系统】运行时间
	data.SystemBootTime = ag.config.systemBootTime   //【系统】开机时间
	data.CpuCores = ag.config.cpuCores               //【CPU】核数
	data.CpuModelName = ag.config.cpuModelName       //【CPU】型号名称
	data.CpuMhz = ag.config.cpuMhz                   //【CPU】兆赫

	err := ag.gormClient.Table(ag.gormConfig.tableName).Create(&data).Error
	if err != nil {
		log.Printf("记录接口日志错误：%s", err)
		log.Printf("记录接口日志数据：%+v", data)
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
