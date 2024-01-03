package golog

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"time"
)

// 结构体
type apiSLog struct {
	TraceID            string                 `json:"trace_id,omitempty"`
	RequestTime        time.Time              `json:"request_time,omitempty"`
	RequestUri         string                 `json:"request_uri,omitempty"`
	RequestUrl         string                 `json:"request_url,omitempty"`
	RequestApi         string                 `json:"request_api,omitempty"`
	RequestMethod      string                 `json:"request_method,omitempty"`
	RequestParams      map[string]interface{} `json:"request_params,omitempty"`
	RequestHeader      map[string]string      `json:"request_header,omitempty"`
	ResponseHeader     map[string][]string    `json:"response_header,omitempty"`
	ResponseStatusCode int                    `json:"response_status_code,omitempty"`
	ResponseBody       map[string]interface{} `json:"response_body,omitempty"`
	ResponseTime       time.Time              `json:"response_time,omitempty,omitempty"`
}

// Middleware 中间件
func (al *ApiSLog) Middleware(ctx context.Context, request gorequest.Response) {
	data := apiSLog{
		TraceID:            gotrace_id.GetTraceIdContext(ctx),
		RequestTime:        request.RequestTime,
		RequestUri:         request.RequestUri,
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,
		RequestApi:         gourl.UriParse(request.RequestUri).Path,
		RequestMethod:      request.RequestMethod,
		RequestParams:      request.RequestParams,
		RequestHeader:      request.RequestHeader,
		ResponseHeader:     request.ResponseHeader,
		ResponseStatusCode: request.ResponseStatusCode,
		ResponseBody:       gojson.JsonDecodeNoError(string(request.ResponseBody)),
		ResponseTime:       request.ResponseTime,
	}
	if al.slog.status {
		al.slog.client.WithTraceId(ctx).Info("Middleware",
			"request_time", data.RequestTime,
			"request_uri", data.RequestUri,
			"request_url", data.RequestUrl,
			"request_api", data.RequestApi,
			"request_method", data.RequestMethod,
			"request_params", data.RequestParams,
			"request_header", data.RequestHeader,
			"response_header", data.ResponseHeader,
			"response_status_code", data.ResponseStatusCode,
			"response_body", data.ResponseBody,
			"response_time", data.ResponseTime,
		)
	}
}

// MiddlewareXml 中间件
func (al *ApiSLog) MiddlewareXml(ctx context.Context, request gorequest.Response) {
	data := apiSLog{
		TraceID:            gotrace_id.GetTraceIdContext(ctx),
		RequestTime:        request.RequestTime,
		RequestUri:         request.RequestUri,
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,
		RequestApi:         gourl.UriParse(request.RequestUri).Path,
		RequestMethod:      request.RequestMethod,
		RequestParams:      request.RequestParams,
		RequestHeader:      request.RequestHeader,
		ResponseHeader:     request.ResponseHeader,
		ResponseStatusCode: request.ResponseStatusCode,
		ResponseBody:       gojson.XmlDecodeNoError(request.ResponseBody),
		ResponseTime:       request.ResponseTime,
	}
	if al.slog.status {
		al.slog.client.WithTraceId(ctx).Info("MiddlewareXml",
			"request_time", data.RequestTime,
			"request_uri", data.RequestUri,
			"request_url", data.RequestUrl,
			"request_api", data.RequestApi,
			"request_method", data.RequestMethod,
			"request_params", data.RequestParams,
			"request_header", data.RequestHeader,
			"response_header", data.ResponseHeader,
			"response_status_code", data.ResponseStatusCode,
			"response_body", data.ResponseBody,
			"response_time", data.ResponseTime,
		)
	}
}

// MiddlewareCustom 中间件
func (al *ApiSLog) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	data := apiSLog{
		TraceID:            gotrace_id.GetTraceIdContext(ctx),
		RequestTime:        request.RequestTime,
		RequestUri:         request.RequestUri,
		RequestUrl:         gourl.UriParse(request.RequestUri).Url,
		RequestApi:         api,
		RequestMethod:      request.RequestMethod,
		RequestParams:      request.RequestParams,
		RequestHeader:      request.RequestHeader,
		ResponseHeader:     request.ResponseHeader,
		ResponseStatusCode: request.ResponseStatusCode,
		ResponseBody:       gojson.JsonDecodeNoError(string(request.ResponseBody)),
		ResponseTime:       request.ResponseTime,
	}
	if al.slog.status {
		al.slog.client.WithTraceId(ctx).Info("MiddlewareCustom",
			"request_time", data.RequestTime,
			"request_uri", data.RequestUri,
			"request_url", data.RequestUrl,
			"request_api", data.RequestApi,
			"request_method", data.RequestMethod,
			"request_params", data.RequestParams,
			"request_header", data.RequestHeader,
			"response_header", data.ResponseHeader,
			"response_status_code", data.ResponseStatusCode,
			"response_body", data.ResponseBody,
			"response_time", data.ResponseTime,
		)
	}
}
