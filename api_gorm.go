package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/gorequest"
)

// ApiGorm 接口日志
type ApiGorm struct {
	gormClient *dorm.GormClient // 数据库驱动
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
	gormConfig struct {
		stats     bool   // 状态
		tableName string // 表名
	}
}

// ApiGormFun 接口日志驱动
type ApiGormFun func() *ApiGorm

// NewApiGorm 创建接口实例化
func NewApiGorm(ctx context.Context, systemOutsideIp string, gormClient *dorm.GormClient, gormTableName string) (*ApiGorm, error) {

	gl := &ApiGorm{}

	// 配置信息
	if systemOutsideIp == "" {
		return nil, errors.New("没有设置外网IP")
	}
	gl.setConfig(ctx, systemOutsideIp)

	if gormClient == nil || gormClient.GetDb() == nil {
		gl.gormConfig.stats = false
	} else {

		gl.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			gl.gormConfig.tableName = gormTableName
		}

		// 创建模型
		gl.gormAutoMigrate(ctx)

		gl.gormConfig.stats = true

	}

	return gl, nil
}

// Middleware 中间件
func (gl *ApiGorm) Middleware(ctx context.Context, request gorequest.Response) {
	if gl.gormConfig.stats {
		gl.gormMiddleware(ctx, request)
	}
}

// MiddlewareXml 中间件
func (gl *ApiGorm) MiddlewareXml(ctx context.Context, request gorequest.Response) {
	if gl.gormConfig.stats {
		gl.gormMiddlewareXml(ctx, request)
	}
}

// MiddlewareCustom 中间件
func (gl *ApiGorm) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	if gl.gormConfig.stats {
		gl.gormMiddlewareCustom(ctx, api, request)
	}
}
