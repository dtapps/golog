package golog

import (
	"github.com/gin-gonic/gin"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gourl"
	"time"
)

// 结构体
type ginSLog struct {
	TraceID            string                 `json:"trace_id,omitempty"`             //【系统】跟踪编号
	RequestTime        time.Time              `json:"request_time,omitempty"`         //【请求】时间
	RequestUri         string                 `json:"request_uri,omitempty"`          //【请求】请求链接 域名+路径+参数
	RequestUrl         string                 `json:"request_url,omitempty"`          //【请求】请求链接 域名+路径
	RequestApi         string                 `json:"request_api,omitempty"`          //【请求】请求接口 路径
	RequestMethod      string                 `json:"request_method,omitempty"`       //【请求】请求方式
	RequestProto       string                 `json:"request_proto,omitempty"`        //【请求】请求协议
	RequestUa          string                 `json:"request_ua,omitempty"`           //【请求】请求UA
	RequestReferer     string                 `json:"request_referer,omitempty"`      //【请求】请求referer
	RequestBody        string                 `json:"request_body,omitempty"`         //【请求】请求主体
	RequestUrlQuery    map[string][]string    `json:"request_url_query,omitempty"`    //【请求】请求URL参数
	RequestIP          string                 `json:"request_ip,omitempty"`           //【请求】请求客户端Ip
	RequestIpCountry   string                 `json:"request_ip_country,omitempty"`   //【请求】请求客户端城市
	RequestIpProvince  string                 `json:"request_ip_province,omitempty"`  //【请求】请求客户端省份
	RequestIpCity      string                 `json:"request_ip_city,omitempty"`      //【请求】请求客户端城市
	RequestIpIsp       string                 `json:"request_ip_isp,omitempty"`       //【请求】请求客户端运营商
	RequestIpLongitude float64                `json:"request_ip_longitude,omitempty"` //【请求】请求客户端经度
	RequestIpLatitude  float64                `json:"request_ip_latitude,omitempty"`  //【请求】请求客户端纬度
	RequestHeader      map[string][]string    `json:"request_header,omitempty"`       //【请求】请求头
	RequestAllContent  map[string]interface{} `json:"request_all_content,omitempty"`  // 【请求】请求全部内容
	ResponseTime       time.Time              `json:"response_time,omitempty"`        //【返回】时间
	ResponseCode       int                    `json:"response_code,omitempty"`        //【返回】状态码
	ResponseMsg        string                 `json:"response_msg,omitempty"`         //【返回】描述
	ResponseData       string                 `json:"response_data,omitempty"`        //【返回】数据
	CostTime           int64                  `json:"cost_time,omitempty"`            //【系统】花费时间
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
		"request_ua", data.RequestUa,
		"request_referer", data.RequestReferer,
		"request_body", data.RequestBody,
		"request_url_query", data.RequestUrlQuery,
		"request_ip", data.RequestIP,
		"request_ip_country", data.RequestIpCountry,
		"request_ip_province", data.RequestIpProvince,
		"request_ip_city", data.RequestIpCity,
		"request_ip_isp", data.RequestIpIsp,
		"request_ip_longitude", data.RequestIpLongitude,
		"request_ip_latitude", data.RequestIpLatitude,
		"request_header", data.RequestHeader,
		"request_all_content", data.RequestAllContent,
		"response_time", data.ResponseTime,
		"response_code", data.ResponseCode,
		"response_msg", data.ResponseMsg,
		"response_data", data.ResponseData,
		"cost_time", data.CostTime,
	)
}

func (gl *GinSLog) recordJson(ginCtx *gin.Context, traceId string, requestTime time.Time, paramsBody gorequest.Params, responseCode int, responseBody string, startTime, endTime int64, ipInfo goip.AnalyseResult) {
	data := ginSLog{
		TraceID:            traceId,                                                      //【系统】跟踪编号
		RequestTime:        requestTime,                                                  //【请求】时间
		RequestUrl:         ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:         gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:      ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:       ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:          ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:     ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:    ginCtx.Request.URL.Query(),                                   //【请求】请求URL参数
		RequestIP:          ipInfo.Ip,                                                    //【请求】请求客户端Ip
		RequestIpCountry:   ipInfo.Country,                                               //【请求】请求客户端城市
		RequestIpProvince:  ipInfo.Province,                                              //【请求】请求客户端省份
		RequestIpCity:      ipInfo.City,                                                  //【请求】请求客户端城市
		RequestIpIsp:       ipInfo.Isp,                                                   //【请求】请求客户端运营商
		RequestIpLatitude:  ipInfo.LocationLatitude,                                      //【请求】请求客户端纬度
		RequestIpLongitude: ipInfo.LocationLongitude,                                     //【请求】请求客户端经度
		RequestHeader:      ginCtx.Request.Header,                                        //【请求】请求头
		RequestAllContent:  paramsBody,                                                   //【请求】请求全部内容
		ResponseTime:       gotime.Current().Time,                                        //【返回】时间
		ResponseCode:       responseCode,                                                 //【返回】状态码
		ResponseData:       responseBody,                                                 //【返回】数据
		CostTime:           endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}
	gl.record("json", data)
}
