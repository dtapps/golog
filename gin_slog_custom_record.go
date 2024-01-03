package golog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
)

// 结构体
type ginSLogCustom struct {
	TraceID         string `json:"trace_id,omitempty"`          //【系统】跟踪编号
	RequestUri      string `json:"request_uri,omitempty"`       //【请求】请求链接 域名+路径+参数
	RequestUrl      string `json:"request_url,omitempty"`       //【请求】请求链接 域名+路径
	RequestApi      string `json:"request_api,omitempty"`       //【请求】请求接口 路径
	RequestMethod   string `json:"request_method,omitempty"`    //【请求】请求方式
	RequestProto    string `json:"request_proto,omitempty"`     //【请求】请求协议
	RequestUa       string `json:"request_ua,omitempty"`        //【请求】请求UA
	RequestReferer  string `json:"request_referer,omitempty"`   //【请求】请求referer
	RequestUrlQuery string `json:"request_url_query,omitempty"` //【请求】请求URL参数
	RequestHeader   string `json:"request_header,omitempty"`    //【请求】请求头
	RequestIP       string `json:"request_ip,omitempty"`        //【请求】请求客户端IP
	CustomID        string `json:"custom_id,omitempty"`         //【日志】自定义编号
	CustomType      string `json:"custom_type,omitempty"`       //【日志】自定义类型
	CustomContent   string `json:"custom_content,omitempty"`    //【日志】自定义内容
}

type GinCustomClientGinRecordOperation struct {
	slogClient *SLog          // 日志服务
	data       *ginSLogCustom // 数据
}

// GinRecord 记录日志
func (c *GinSLogCustom) GinRecord(ginCtx *gin.Context) *GinCustomClientGinRecordOperation {
	operation := &GinCustomClientGinRecordOperation{
		slogClient: c.slog.client,
	}
	operation.data = new(ginSLogCustom)
	operation.data.TraceID = gotrace_id.GetGinTraceId(ginCtx) // 【系统】跟踪编号
	if ginCtx.Request.TLS == nil {
		operation.data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		operation.data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}
	operation.data.RequestUrl = ginCtx.Request.RequestURI                                    //【请求】请求链接 域名+路径
	operation.data.RequestApi = gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI) //【请求】请求接口 路径
	operation.data.RequestMethod = ginCtx.Request.Method                                     //【请求】请求方式
	operation.data.RequestProto = ginCtx.Request.Proto                                       //【请求】请求协议
	operation.data.RequestUa = ginCtx.Request.UserAgent()                                    //【请求】请求UA
	operation.data.RequestReferer = ginCtx.Request.Referer()                                 //【请求】请求referer
	operation.data.RequestUrlQuery = gojson.JsonEncodeNoError(ginCtx.Request.URL.Query())    //【请求】请求URL参数
	operation.data.RequestHeader = gojson.JsonEncodeNoError(ginCtx.Request.Header)           //【请求】请求头
	operation.data.RequestIP = gorequest.ClientIp(ginCtx.Request)                            //【请求】请求客户端Ip
	return operation
}

func (o *GinCustomClientGinRecordOperation) CustomInfo(customId any, customType any, customContent any) *GinCustomClientGinRecordOperation {
	o.data.CustomID = fmt.Sprintf("%s", customId)           //【日志】自定义编号
	o.data.CustomType = fmt.Sprintf("%s", customType)       //【日志】自定义类型
	o.data.CustomContent = fmt.Sprintf("%s", customContent) //【日志】自定义内容
	return o
}

func (o *GinCustomClientGinRecordOperation) CreateData() {
	o.slogClient.WithTraceIdStr(o.data.TraceID).Info("custom",
		"request_uri", o.data.RequestUri,
		"request_url", o.data.RequestUrl,
		"request_api", o.data.RequestApi,
		"request_method", o.data.RequestMethod,
		"request_proto", o.data.RequestProto,
		"request_ua", o.data.RequestUa,
		"request_referer", o.data.RequestReferer,
		"request_url_query", o.data.RequestUrlQuery,
		"request_header", o.data.RequestHeader,
		"request_ip", o.data.RequestIP,
		"custom_id", o.data.CustomID,
		"custom_type", o.data.CustomType,
		"custom_content", o.data.CustomContent,
	)
}

func (o *GinCustomClientGinRecordOperation) CreateDataNoError() {
	o.slogClient.WithTraceIdStr(o.data.TraceID).Info("custom",
		"request_uri", o.data.RequestUri,
		"request_url", o.data.RequestUrl,
		"request_api", o.data.RequestApi,
		"request_method", o.data.RequestMethod,
		"request_proto", o.data.RequestProto,
		"request_ua", o.data.RequestUa,
		"request_referer", o.data.RequestReferer,
		"request_url_query", o.data.RequestUrlQuery,
		"request_header", o.data.RequestHeader,
		"request_ip", o.data.RequestIP,
		"custom_id", o.data.CustomID,
		"custom_type", o.data.CustomType,
		"custom_content", o.data.CustomContent,
	)
}
