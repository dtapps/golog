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

	ag.config.systemInsideIp = gorequest.GetInsideIp(ctx)
	ag.config.systemOutsideIp = systemOutsideIp

	ag.config.goVersion = runtime.Version()
	ag.config.sdkVersion = Version

}
