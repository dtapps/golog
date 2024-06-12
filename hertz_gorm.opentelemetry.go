package golog

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetTrace 设置OpenTelemetry链路追踪
func (hg *HertzGorm) SetTrace(trace bool) {
	hg.trace = trace
}

// TraceStartSpan 开始OpenTelemetry链路追踪状态
func (hg *HertzGorm) TraceStartSpan(ctx context.Context) context.Context {
	if hg.trace {
		tr := otel.Tracer("go.dtapp.net/golog", trace.WithInstrumentationVersion(Version))
		ctx, hg.span = tr.Start(ctx, "golog.hertz")
	}
	return ctx
}

// TraceEndSpan 结束OpenTelemetry链路追踪状态
func (hg *HertzGorm) TraceEndSpan() {
	if hg.trace {
		hg.span.End()
	}
}

// TraceSetAttributes 设置OpenTelemetry链路追踪属性
func (hg *HertzGorm) TraceSetAttributes(kv ...attribute.KeyValue) {
	if hg.trace {
		hg.span.SetAttributes(kv...)
	}
}

// TraceSetStatus 设置OpenTelemetry链路追踪状态
func (hg *HertzGorm) TraceSetStatus(code codes.Code, description string) {
	if hg.trace {
		hg.span.SetStatus(code, description)
	}
}

// TraceGetTraceID 获取OpenTelemetry链路追踪TraceID
func (hg *HertzGorm) TraceGetTraceID() string {
	if hg.trace {
		return hg.span.SpanContext().TraceID().String()
	}
	return ""
}

// TraceGetSpanID 获取OpenTelemetry链路追踪SpanID
func (hg *HertzGorm) TraceGetSpanID() string {
	if hg.trace {
		return hg.span.SpanContext().SpanID().String()
	}
	return ""
}
