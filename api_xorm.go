package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gorequest"
	"xorm.io/xorm"
)

// ApiXorm 接口日志
type ApiXorm struct {
	xormClient *xorm.Engine // 数据库驱动
	config     struct {
		systemHostname      string  // 主机名
		systemOs            string  // 系统类型
		systemVersion       string  // 系统版本
		systemKernel        string  // 系统内核
		systemKernelVersion string  // 系统内核版本
		systemBootTime      uint64  // 系统开机时间
		cpuCores            int     // CPU核数
		cpuModelName        string  // CPU型号名称
		cpuMhz              float64 // CPU兆赫
		systemInsideIp      string  // 内网ip
		systemOutsideIp     string  // 外网ip
		goVersion           string  // go版本
		sdkVersion          string  // sdk版本
	}
	xormConfig struct {
		stats     bool   // 状态
		tableName string // 表名
	}
}

// ApiXormFun 接口日志驱动
type ApiXormFun func() *ApiXorm

// NewApiXorm 创建接口实例化
func NewApiXorm(ctx context.Context, systemOutsideIp string, xormClient *xorm.Engine, xormTableName string) (*ApiXorm, error) {

	gl := &ApiXorm{}

	// 配置信息
	if systemOutsideIp == "" {
		return nil, errors.New("没有设置外网IP")
	}
	gl.setConfig(ctx, systemOutsideIp)

	if xormClient == nil {
		gl.xormConfig.stats = false
	} else {

		gl.xormClient = xormClient

		if xormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			gl.xormConfig.tableName = xormTableName
		}

		// 创建模型
		gl.xormSync(ctx)

		gl.xormConfig.stats = true

	}

	return gl, nil
}

// Middleware 中间件
func (ag *ApiXorm) Middleware(ctx context.Context, request gorequest.Response) {
	if ag.xormConfig.stats {
		ag.xormMiddleware(ctx, request)
	}
}

// MiddlewareXml 中间件
func (ag *ApiXorm) MiddlewareXml(ctx context.Context, request gorequest.Response) {
	if ag.xormConfig.stats {
		ag.xormMiddlewareXml(ctx, request)
	}
}

// MiddlewareCustom 中间件
func (ag *ApiXorm) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	if ag.xormConfig.stats {
		ag.xormMiddlewareCustom(ctx, api, request)
	}
}
