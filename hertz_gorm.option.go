package golog

// SetTrace 设置OpenTelemetry链路追踪
func (hg *HertzGorm) SetTrace(trace bool) {
	hg.trace = trace
}
