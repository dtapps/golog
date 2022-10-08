package golog

import (
	"context"
	"go.dtapp.net/goip"
	"go.mongodb.org/mongo-driver/version"
	"runtime"
)

func (c *ApiClient) setConfig(ctx context.Context) {

	info := getSystem()

	c.config.systemHostname = info.SystemHostname
	c.config.systemOs = info.SystemOs
	c.config.systemVersion = info.SystemVersion
	c.config.systemKernel = info.SystemKernel
	c.config.systemKernelVersion = info.SystemKernelVersion
	c.config.systemBootTime = info.SystemBootTime
	c.config.cpuCores = info.CpuCores
	c.config.cpuModelName = info.CpuModelName
	c.config.cpuMhz = info.CpuMhz

	c.config.systemInsideIp = goip.GetInsideIp(ctx)

	c.config.sdkVersion = Version
	c.config.goVersion = runtime.Version()

	c.config.mongoSdkVersion = version.Driver
}
