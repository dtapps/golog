package golog

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/requestid"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gotime"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
)

// HertzLogFunc Hertz框架日志函数
type HertzLogFunc func(ctx context.Context, response *GormHertzLogModel)

// HertzGorm 框架日志
type HertzGorm struct {
	hertzLogFunc HertzLogFunc // Hertz框架日志函数
	trace        bool         // OpenTelemetry链路追踪
	span         trace.Span   // OpenTelemetry链路追踪
}

// HertzGormFun *HertzGorm 框架日志驱动
type HertzGormFun func() *HertzGorm

// NewHertzGorm 创建框架实例化
func NewHertzGorm(ctx context.Context) (*HertzGorm, error) {
	hg := &HertzGorm{}

	hg.trace = true
	return hg, nil
}

// Middleware 中间件
func (hg *HertzGorm) Middleware() app.HandlerFunc {
	return func(ctx context.Context, h *app.RequestContext) {

		// OpenTelemetry链路追踪
		ctx = hg.TraceStartSpan(ctx)
		defer hg.TraceEndSpan()

		// 开始时间
		start := time.Now().UTC()

		// 模型
		var log = GormHertzLogModel{}

		// 请求时间
		log.RequestTime = gotime.Current().Time

		// 处理请求
		h.Next(ctx)

		// 结束时间
		end := time.Now().UTC()

		// 请求消耗时长
		log.RequestCostTime = end.Sub(start).Milliseconds()

		// 响应时间
		log.ResponseTime = gotime.Current().Time

		// 输出路由日志
		hlog.CtxTracef(ctx, "status=%d cost=%d method=%s full_path=%s client_ip=%s host=%s",
			h.Response.StatusCode(),
			log.RequestCostTime,
			h.Request.Header.Method(),
			h.Request.URI().PathOriginal(),
			h.ClientIP(),
			h.Request.Host(),
		)

		// 跟踪编号
		log.TraceID = hg.TraceGetTraceID()

		// 请求编号
		log.RequestID = requestid.Get(h)

		// 请求主机
		log.RequestHost = string(h.Request.Host())

		// 请求地址
		log.RequestPath = string(h.Request.URI().Path())

		// 请求参数
		log.RequestQuery = gojson.JsonEncodeNoError(gojson.ParseQueryString(string(h.Request.QueryString())))

		// 请求方式
		log.RequestMethod = string(h.Request.Header.Method())

		// 请求协议
		log.RequestScheme = string(h.Request.Scheme())

		// 请求类型
		log.RequestContentType = string(h.ContentType())

		if strings.Contains(log.RequestContentType, consts.MIMEApplicationHTMLForm) {
			log.RequestBody = gojson.JsonEncodeNoError(gojson.ParseQueryString(string(h.Request.Body())))
		} else if strings.Contains(log.RequestContentType, consts.MIMEMultipartPOSTForm) {
			log.RequestBody = string(h.Request.Body())
		} else if strings.Contains(log.RequestContentType, consts.MIMEApplicationJSON) {
			log.RequestBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(h.Request.Body())))
		} else {
			log.RequestBody = string(h.Request.Body())
		}

		// 请求IP
		log.RequestClientIP = h.ClientIP()

		// 请求UA
		log.RequestUserAgent = string(h.UserAgent())

		// 请求头
		requestHeader := make(map[string][]string)
		h.Request.Header.VisitAll(func(k, v []byte) {
			requestHeader[string(k)] = append(requestHeader[string(k)], string(v))
		})
		log.RequestHeader = gojson.JsonEncodeNoError(requestHeader)

		// 响应头
		responseHeader := make(map[string][]string)
		h.Response.Header.VisitAll(func(k, v []byte) {
			responseHeader[string(k)] = append(responseHeader[string(k)], string(v))
		})
		log.ResponseHeader = gojson.JsonEncodeNoError(responseHeader)

		// 响应状态
		log.ResponseStatusCode = h.Response.StatusCode()

		// 响应内容
		if gojson.IsValidJSON(string(h.Response.Body())) {
			log.ResponseBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(h.Response.Body())))
		} else {
			log.ResponseBody = string(h.Response.Body())
		}

		// OpenTelemetry链路追踪
		hg.TraceSetAttributes(attribute.String("request.id", log.RequestID))
		hg.TraceSetAttributes(attribute.String("request.time", log.RequestTime.Format(gotime.DateTimeFormat)))
		hg.TraceSetAttributes(attribute.String("request.host", log.RequestHost))
		hg.TraceSetAttributes(attribute.String("request.path", log.RequestPath))
		hg.TraceSetAttributes(attribute.String("request.query", log.RequestQuery))
		hg.TraceSetAttributes(attribute.String("request.method", log.RequestMethod))
		hg.TraceSetAttributes(attribute.String("request.scheme", log.RequestScheme))
		hg.TraceSetAttributes(attribute.String("request.content_type", log.RequestContentType))
		hg.TraceSetAttributes(attribute.String("request.body", log.RequestBody))
		hg.TraceSetAttributes(attribute.String("request.client_ip", log.RequestClientIP))
		hg.TraceSetAttributes(attribute.String("request.user_agent", log.RequestClientIP))
		hg.TraceSetAttributes(attribute.String("request.header", log.RequestHeader))
		hg.TraceSetAttributes(attribute.Int64("request.cost_time", log.RequestCostTime))
		hg.TraceSetAttributes(attribute.String("response.time", log.ResponseTime.Format(gotime.DateTimeFormat)))
		hg.TraceSetAttributes(attribute.String("response.header", log.ResponseHeader))
		hg.TraceSetAttributes(attribute.Int("response.status_code", log.ResponseStatusCode))
		hg.TraceSetAttributes(attribute.String("response.body", log.ResponseBody))

		// 调用Hertz框架日志函数
		if hg.hertzLogFunc != nil {
			hg.hertzLogFunc(ctx, &log)
		}

	}
}
