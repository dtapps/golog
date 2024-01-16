package golog

import (
	"context"
	"go.dtapp.net/gorequest"
	"runtime"
)

func (gm *GinMongo) setConfig(ctx context.Context, systemOutsideIp string) {

	info := getSystem()

	gm.config.systemHostname = info.SystemHostname
	gm.config.systemOs = info.SystemOs
	gm.config.systemVersion = info.SystemVersion
	gm.config.systemKernel = info.SystemKernel
	gm.config.systemKernelVersion = info.SystemKernelVersion
	gm.config.systemUpTime = info.SystemUpTime
	gm.config.systemBootTime = info.SystemBootTime
	gm.config.cpuCores = info.CpuCores
	gm.config.cpuModelName = info.CpuModelName
	gm.config.cpuMhz = info.CpuMhz

	gm.config.systemInsideIP = gorequest.GetInsideIp(ctx)
	gm.config.systemOutsideIP = systemOutsideIp

	gm.config.sdkVersion = Version
	gm.config.goVersion = runtime.Version()

}

// ConfigSLogClientFun 日志配置
func (gm *GinMongo) ConfigSLogClientFun(sLogFun SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		gm.slog.client = sLog
		gm.slog.status = true
	}
}
