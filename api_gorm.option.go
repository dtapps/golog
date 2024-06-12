package golog

// SetTrace 设置OpenTelemetry链路追踪
func (ag *ApiGorm) SetTrace(trace bool) {
	ag.trace = trace
}
