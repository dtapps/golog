package golog

// ConfigSLogClientFun 日志配置
func (c *GinSLog) ConfigSLogClientFun(sLogFun SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		c.slog.client = sLog
		c.slog.status = true
	}
}

// ConfigSLogResultClientFun 日志配置然后返回
func (c *GinSLog) ConfigSLogResultClientFun(sLogFun SLogFun) *GinSLog {
	sLog := sLogFun()
	if sLog != nil {
		c.slog.client = sLog
		c.slog.status = true
	}
	return c
}
