package golog

import (
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"time"
)

// 结构体
type ginSLog struct {
	TraceID       string                 `json:"trace_id,omitempty"`       //【系统】跟踪编号
	RequestTime   time.Time              `json:"request_time,omitempty"`   //【请求】时间
	RequestUri    string                 `json:"request_uri,omitempty"`    //【请求】请求链接 域名+路径+参数
	RequestUrl    string                 `json:"request_url,omitempty"`    //【请求】请求链接 域名+路径
	RequestApi    string                 `json:"request_api,omitempty"`    //【请求】请求接口 路径
	RequestMethod string                 `json:"request_method,omitempty"` //【请求】请求方式
	RequestProto  string                 `json:"request_proto,omitempty"`  //【请求】请求协议
	RequestBody   map[string]interface{} `json:"request_body,omitempty"`   //【请求】请求参数
	RequestIP     string                 `json:"request_ip,omitempty"`     //【请求】请求客户端IP
	RequestHeader map[string][]string    `json:"request_header,omitempty"` //【请求】请求头
	ResponseTime  time.Time              `json:"response_time,omitempty"`  //【返回】时间
	ResponseCode  int                    `json:"response_code,omitempty"`  //【返回】状态码
	ResponseData  string                 `json:"response_data,omitempty"`  //【返回】数据
	CostTime      int64                  `json:"cost_time,omitempty"`      //【系统】花费时间
}

// record 记录日志
func (gl *GinSLog) record(msg string, data ginSLog) {
	gl.slog.client.WithTraceIdStr(data.TraceID).Info(msg,
		"request_time", data.RequestTime,
		"request_uri", data.RequestUri,
		"request_url", data.RequestUrl,
		"request_api", data.RequestApi,
		"request_method", data.RequestMethod,
		"request_proto", data.RequestProto,
		"request_body", data.RequestBody,
		"request_ip", data.RequestIP,
		"request_header", data.RequestHeader,
		"response_time", data.ResponseTime,
		"response_code", data.ResponseCode,
		"response_data", data.ResponseData,
		"cost_time", data.CostTime,
	)
}

func (gl *GinSLog) recordJson(ginCtx *gin.Context, requestTime time.Time, request_body gorequest.Params, responseTime time.Time, responseCode int, responseBody string, costTime int64, requestIp string) {
	data := ginSLog{
		TraceID:       gotrace_id.GetGinTraceId(ginCtx),                             //【系统】跟踪编号
		RequestTime:   requestTime,                                                  //【请求】时间
		RequestUrl:    ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:    gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod: ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:  ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestBody:   request_body,                                                 //【请求】请求参数
		RequestIP:     requestIp,                                                    //【请求】请求客户端IP
		RequestHeader: ginCtx.Request.Header,                                        //【请求】请求头
		ResponseTime:  responseTime,                                                 //【返回】时间
		ResponseCode:  responseCode,                                                 //【返回】状态码
		ResponseData:  responseBody,                                                 //【返回】数据
		CostTime:      costTime,                                                     //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}
	gl.record("json", data)
}
