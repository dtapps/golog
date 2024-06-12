package golog

// SetTrace 设置OpenTelemetry链路追踪
func (gg *GinGorm) SetTrace(trace bool) {
	gg.trace = trace
}
