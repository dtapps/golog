package golog

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// SetTrace 设置OpenTelemetry链路追踪
// TODO: 等待完全删除
func (gg *GinGorm) SetTrace(trace bool) {

}

// TraceStartSpan 开始OpenTelemetry链路追踪状态
func (gg *GinGorm) TraceStartSpan(ctx context.Context) (context.Context, trace.Span) {
	tr := otel.Tracer("go.dtapp.net/golog", trace.WithInstrumentationVersion(Version))
	ctx, span := tr.Start(ctx, "golog.gin", trace.WithSpanKind(trace.SpanKindInternal))
	return ctx, span
}
