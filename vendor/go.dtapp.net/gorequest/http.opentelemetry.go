package gorequest

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetTrace 设置OpenTelemetry链路追踪
func (c *App) SetTrace(trace bool) {
	c.trace = trace
}

// TraceStartSpan 开始OpenTelemetry链路追踪状态
func (c *App) TraceStartSpan(ctx context.Context, spanName string) context.Context {
	if c.trace {
		tr := otel.Tracer("go.dtapp.net/gorequest", trace.WithInstrumentationVersion(Version))
		ctx, c.span = tr.Start(ctx, "gorequest."+spanName)
	}
	return ctx
}

// TraceEndSpan 结束OpenTelemetry链路追踪状态
func (c *App) TraceEndSpan() {
	if c.trace {
		c.span.End()
	}
}

// TraceSetAttributes 设置OpenTelemetry链路追踪属性
func (c *App) TraceSetAttributes(kv ...attribute.KeyValue) {
	if c.trace {
		c.span.SetAttributes(kv...)
	}
}

// TraceSetStatus 设置OpenTelemetry链路追踪状态
func (c *App) TraceSetStatus(code codes.Code, description string) {
	if c.trace {
		c.span.SetStatus(code, description)
	}
}

// TraceGetTraceID 获取OpenTelemetry链路追踪TraceID
func (c *App) TraceGetTraceID() string {
	if c.trace {
		return c.span.SpanContext().TraceID().String()
	}
	return ""
}

// TraceGetSpanID 获取OpenTelemetry链路追踪SpanID
func (c *App) TraceGetSpanID() string {
	if c.trace {
		return c.span.SpanContext().SpanID().String()
	}
	return ""
}
