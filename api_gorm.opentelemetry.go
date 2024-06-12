package golog

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetTrace 设置OpenTelemetry链路追踪
func (ag *ApiGorm) SetTrace(trace bool) {
	ag.trace = trace
}

// TraceStartSpan 开始OpenTelemetry链路追踪状态
func (ag *ApiGorm) TraceStartSpan(ctx context.Context) context.Context {
	if ag.trace {
		tr := otel.Tracer("go.dtapp.net/golog", trace.WithInstrumentationVersion(Version))
		ctx, ag.span = tr.Start(ctx, "golog.api")
	}
	return ctx
}

// TraceEndSpan 结束OpenTelemetry链路追踪状态
func (ag *ApiGorm) TraceEndSpan() {
	if ag.trace {
		ag.span.End()
	}
}

// TraceSetAttributes 设置OpenTelemetry链路追踪属性
func (ag *ApiGorm) TraceSetAttributes(kv ...attribute.KeyValue) {
	if ag.trace {
		ag.span.SetAttributes(kv...)
	}
}

// TraceSetStatus 设置OpenTelemetry链路追踪状态
func (ag *ApiGorm) TraceSetStatus(code codes.Code, description string) {
	if ag.trace {
		ag.span.SetStatus(code, description)
	}
}

// TraceGetTraceID 获取OpenTelemetry链路追踪TraceID
func (ag *ApiGorm) TraceGetTraceID() string {
	if ag.trace {
		return ag.span.SpanContext().TraceID().String()
	}
	return ""
}

// TraceGetSpanID 获取OpenTelemetry链路追踪SpanID
func (ag *ApiGorm) TraceGetSpanID() string {
	if ag.trace {
		return ag.span.SpanContext().SpanID().String()
	}
	return ""
}
