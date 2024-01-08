package golog

import (
	"context"
	"go.dtapp.net/gorequest"
	"runtime"
)

func (gg *GinXorm) setConfig(ctx context.Context, systemOutsideIp string) {

	info := getSystem()

	gg.config.systemHostname = info.SystemHostname
	gg.config.systemOs = info.SystemOs
	gg.config.systemVersion = info.SystemVersion
	gg.config.systemKernel = info.SystemKernel
	gg.config.systemKernelVersion = info.SystemKernelVersion
	gg.config.systemUpTime = info.SystemUpTime
	gg.config.systemBootTime = info.SystemBootTime
	gg.config.cpuCores = info.CpuCores
	gg.config.cpuModelName = info.CpuModelName
	gg.config.cpuMhz = info.CpuMhz

	gg.config.systemInsideIp = gorequest.GetInsideIp(ctx)
	gg.config.systemOutsideIp = systemOutsideIp

	gg.config.sdkVersion = Version
	gg.config.goVersion = runtime.Version()

}
