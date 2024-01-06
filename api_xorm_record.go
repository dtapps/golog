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
type apiXormLog struct {
	LogID              int64     `xorm:"pk;comment('【记录】编号')" json:"log_id,omitempty"`               //【记录】编号
	TraceID            string    `xorm:"index;comment('【系统】跟踪编号')" json:"trace_id,omitempty"`        //【系统】跟踪编号
	RequestTime        time.Time `xorm:"index;comment('【请求】时间')" json:"request_time,omitempty"`      //【请求】时间
	RequestUri         string    `xorm:"comment('【请求】链接')" json:"request_uri,omitempty"`             //【请求】链接
	RequestUrl         string    `xorm:"comment('【请求】链接')" json:"request_url,omitempty"`             //【请求】链接
	RequestApi         string    `xorm:"index;comment('【请求】接口')" json:"request_api,omitempty"`       //【请求】接口
	RequestMethod      string    `xorm:"index;comment('【请求】方式')" json:"request_method,omitempty"`    //【请求】方式
	RequestParams      string    `xorm:"comment('【请求】参数')" json:"request_params,omitempty"`          //【请求】参数
	RequestHeader      string    `xorm:"comment('【请求】头部')" json:"request_header,omitempty"`          //【请求】头部
	RequestIp          string    `xorm:"comment('【请求】请求IP')" json:"request_ip,omitempty"`            //【请求】请求IP
	ResponseHeader     string    `xorm:"comment('【返回】头部')" json:"response_header,omitempty"`         //【返回】头部
	ResponseStatusCode int       `xorm:"comment('【返回】状态码')" json:"response_status_code,omitempty"`   //【返回】状态码
	ResponseBody       string    `xorm:"comment('【返回】数据')" json:"response_body,omitempty"`           //【返回】数据
	ResponseTime       time.Time `xorm:"index;comment('【返回】时间')" json:"response_time,omitempty"`     //【返回】时间
	SystemHostName     string    `xorm:"index;comment('【系统】主机名')" json:"system_host_name,omitempty"` //【系统】主机名
	SystemInsideIp     string    `xorm:"comment('【系统】内网IP')" json:"system_inside_ip,omitempty"`      //【系统】内网IP
	SystemOs           string    `xorm:"index;comment('【系统】系统类型')" json:"system_os,omitempty"`       //【系统】系统类型
	SystemArch         string    `xorm:"index;comment('【系统】系统架构')" json:"system_arch,omitempty"`     //【系统】系统架构
	GoVersion          string    `xorm:"index;comment('【程序】Go版本')" json:"go_version,omitempty"`      //【程序】Go版本
	SdkVersion         string    `xorm:"index;comment('【程序】Sdk版本')" json:"sdk_version,omitempty"`    //【程序】Sdk版本
}

// 创建模型
func (ag *ApiXorm) xormSync(ctx context.Context) {
	if ag.xormConfig.stats == false {
		return
	}

	err := ag.xormClient.Table(ag.xormConfig.tableName).Sync(&apiXormLog{})
	if err != nil {
		log.Printf("创建模型：%s", err)
	}

}

// 记录日志
func (ag *ApiXorm) xormRecord(ctx context.Context, data apiXormLog) {
	if ag.xormConfig.stats == false {
		return
	}

	if utf8.ValidString(data.ResponseBody) == false {
		data.ResponseBody = ""
	}

	data.SystemHostName = ag.config.systemHostname   //【系统】主机名
	data.SystemInsideIp = ag.config.systemInsideIp   //【系统】内网ip
	data.GoVersion = ag.config.goVersion             //【程序】Go版本
	data.SdkVersion = ag.config.systemVersion        //【程序】Sdk版本
	data.TraceID = gotrace_id.GetTraceIdContext(ctx) //【记录】跟踪编号
	data.RequestIp = ag.config.systemOutsideIp       //【请求】请求Ip
	data.SystemOs = ag.config.systemOs               //【系统】系统类型
	data.SystemArch = ag.config.systemKernel         //【系统】系统架构

	_, err := ag.xormClient.Table(ag.xormConfig.tableName).Insert(&data)
	if err != nil {
		log.Printf("记录接口日志错误：%s", err)
		log.Printf("记录接口日志数据：%+v", data)
	}
}

// 中间件
func (ag *ApiXorm) xormMiddleware(ctx context.Context, request gorequest.Response) {
	data := apiXormLog{
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

	ag.xormRecord(ctx, data)
}

// 中间件
func (ag *ApiXorm) xormMiddlewareXml(ctx context.Context, request gorequest.Response) {
	data := apiXormLog{
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

	ag.xormRecord(ctx, data)
}

// 中间件
func (ag *ApiXorm) xormMiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	data := apiXormLog{
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

	ag.xormRecord(ctx, data)
}
