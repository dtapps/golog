package golog

import (
	"context"
	"go.dtapp.net/gorequest"
	"runtime"
)

func (gl *ApiGorm) setConfig(ctx context.Context, systemOutsideIp string) {

	info := getSystem()

	gl.config.systemHostname = info.SystemHostname
	gl.config.systemOs = info.SystemOs
	gl.config.systemVersion = info.SystemVersion
	gl.config.systemKernel = info.SystemKernel
	gl.config.systemKernelVersion = info.SystemKernelVersion
	gl.config.systemBootTime = info.SystemBootTime
	gl.config.cpuCores = info.CpuCores
	gl.config.cpuModelName = info.CpuModelName
	gl.config.cpuMhz = info.CpuMhz

	gl.config.systemInsideIp = gorequest.GetInsideIp(ctx)
	gl.config.systemOutsideIp = systemOutsideIp

	gl.config.goVersion = runtime.Version()
	gl.config.sdkVersion = Version

}
