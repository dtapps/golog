package golog

import (
	"context"
	"go.dtapp.net/gorequest"
	"runtime"
)

func (ag *ApiGorm) setConfig(ctx context.Context, systemOutsideIp string) {

	info := getSystem()

	ag.config.systemHostname = info.SystemHostname
	ag.config.systemOs = info.SystemOs
	ag.config.systemVersion = info.SystemVersion
	ag.config.systemKernel = info.SystemKernel
	ag.config.systemKernelVersion = info.SystemKernelVersion
	ag.config.systemUpTime = info.SystemUpTime
	ag.config.systemBootTime = info.SystemBootTime
	ag.config.cpuCores = info.CpuCores
	ag.config.cpuModelName = info.CpuModelName
	ag.config.cpuMhz = info.CpuMhz

	ag.config.systemInsideIP = gorequest.GetInsideIp(ctx)
	ag.config.systemOutsideIP = systemOutsideIp

	ag.config.goVersion = runtime.Version()
	ag.config.sdkVersion = Version

}

// ConfigSLogClientFun 日志配置
func (ag *ApiGorm) ConfigSLogClientFun(sLogFun SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		ag.slog.client = sLog
		ag.slog.status = true
	}
}
