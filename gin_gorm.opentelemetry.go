package golog

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetTrace 设置OpenTelemetry链路追踪
func (gg *GinGorm) SetTrace(trace bool) {
	gg.trace = trace
}

// TraceStartSpan 开始OpenTelemetry链路追踪状态
func (gg *GinGorm) TraceStartSpan(ctx context.Context) context.Context {
	if gg.trace {
		tr := otel.Tracer("go.dtapp.net/golog", trace.WithInstrumentationVersion(Version))
		ctx, gg.span = tr.Start(ctx, "golog.gin")
	}
	return ctx
}

// TraceEndSpan 结束OpenTelemetry链路追踪状态
func (gg *GinGorm) TraceEndSpan() {
	if gg.trace {
		gg.span.End()
	}
}

// TraceSetAttributes 设置OpenTelemetry链路追踪属性
func (gg *GinGorm) TraceSetAttributes(kv ...attribute.KeyValue) {
	if gg.trace {
		gg.span.SetAttributes(kv...)
	}
}

// TraceSetStatus 设置OpenTelemetry链路追踪状态
func (gg *GinGorm) TraceSetStatus(code codes.Code, description string) {
	if gg.trace {
		gg.span.SetStatus(code, description)
	}
}

// TraceGetTraceID 获取OpenTelemetry链路追踪TraceID
func (gg *GinGorm) TraceGetTraceID() string {
	if gg.trace {
		return gg.span.SpanContext().TraceID().String()
	}
	return ""
}

// TraceGetSpanID 获取OpenTelemetry链路追踪SpanID
func (gg *GinGorm) TraceGetSpanID() string {
	if gg.trace {
		return gg.span.SpanContext().SpanID().String()
	}
	return ""
}
