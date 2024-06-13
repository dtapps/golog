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
		ctx, hg.span = tr.Start(ctx, "golog.hertz", trace.WithSpanKind(trace.SpanKindInternal))
	}
	return ctx
}

// TraceEndSpan 结束OpenTelemetry链路追踪状态
func (hg *HertzGorm) TraceEndSpan() {
	if hg.trace && hg.span != nil {
		hg.span.End()
	}
}

// TraceSetAttributes 设置OpenTelemetry链路追踪属性
func (hg *HertzGorm) TraceSetAttributes(kv ...attribute.KeyValue) {
	if hg.trace && hg.span != nil {
		hg.span.SetAttributes(kv...)
	}
}

// TraceSetStatus 设置OpenTelemetry链路追踪状态
func (hg *HertzGorm) TraceSetStatus(code codes.Code, description string) {
	if hg.trace && hg.span != nil {
		hg.span.SetStatus(code, description)
	}
}

// TraceRecordError 记录OpenTelemetry链路追踪错误
func (hg *HertzGorm) TraceRecordError(err error, options ...trace.EventOption) {
	if hg.trace && hg.span != nil {
		hg.span.RecordError(err, options...)
	}
}

// TraceGetTraceID 获取OpenTelemetry链路追踪TraceID
func (hg *HertzGorm) TraceGetTraceID() (traceID string) {
	if hg.trace && hg.span != nil {
		traceID = hg.span.SpanContext().TraceID().String()
	}
	return traceID
}

// TraceGetSpanID 获取OpenTelemetry链路追踪SpanID
func (hg *HertzGorm) TraceGetSpanID() (spanID string) {
	if hg.trace && hg.span != nil {
		spanID = hg.span.SpanContext().SpanID().String()
	}
	return spanID
}
