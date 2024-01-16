package golog

import (
	"context"
	"go.dtapp.net/gorequest"
	"runtime"
)

func (am *ApiMongo) setConfig(ctx context.Context, systemOutsideIp string) {

	info := getSystem()

	am.config.systemHostname = info.SystemHostname
	am.config.systemOs = info.SystemOs
	am.config.systemVersion = info.SystemVersion
	am.config.systemKernel = info.SystemKernel
	am.config.systemKernelVersion = info.SystemKernelVersion
	am.config.systemUpTime = info.SystemUpTime
	am.config.systemBootTime = info.SystemBootTime
	am.config.cpuCores = info.CpuCores
	am.config.cpuModelName = info.CpuModelName
	am.config.cpuMhz = info.CpuMhz

	am.config.systemInsideIP = gorequest.GetInsideIp(ctx)
	am.config.systemOutsideIP = systemOutsideIp

	am.config.goVersion = runtime.Version()
	am.config.sdkVersion = Version

}

// ConfigSLogClientFun 日志配置
func (am *ApiMongo) ConfigSLogClientFun(sLogFun SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		am.slog.client = sLog
		am.slog.status = true
	}
}
