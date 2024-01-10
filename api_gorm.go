package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gorequest"
	"gorm.io/gorm"
)

// ApiGorm 接口日志
type ApiGorm struct {
	gormClient *gorm.DB // 数据库驱动
	config     struct {
		systemHostname      string  // 主机名
		systemOs            string  // 系统类型
		systemVersion       string  // 系统版本
		systemKernel        string  // 系统内核
		systemKernelVersion string  // 系统内核版本
		systemUpTime        uint64  // 系统运行时间
		systemBootTime      uint64  // 系统开机时间
		cpuCores            int     // CPU核数
		cpuModelName        string  // CPU型号名称
		cpuMhz              float64 // CPU兆赫
		systemInsideIP      string  // 内网IP
		systemOutsideIP     string  // 外网IP
		goVersion           string  // go版本
		sdkVersion          string  // sdk版本
	}
	gormConfig struct {
		stats     bool   // 状态
		tableName string // 表名
	}
	slog struct {
		status bool  // 状态
		client *SLog // 日志服务
	}
}

// ApiGormFun 接口日志驱动
type ApiGormFun func() *ApiGorm

// NewApiGorm 创建接口实例化
func NewApiGorm(ctx context.Context, systemOutsideIp string, gormClient *gorm.DB, gormTableName string) (*ApiGorm, error) {

	gl := &ApiGorm{}

	// 配置信息
	if systemOutsideIp == "" {
		return nil, errors.New("没有设置外网IP")
	}
	gl.setConfig(ctx, systemOutsideIp)

	if gormClient == nil {
		gl.gormConfig.stats = false
	} else {

		gl.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			gl.gormConfig.tableName = gormTableName
		}

		gl.gormConfig.stats = true

		// 创建模型
		gl.gormAutoMigrate(ctx)

	}

	return gl, nil
}

// Middleware 中间件
func (ag *ApiGorm) Middleware(ctx context.Context, request gorequest.Response) {
	if ag.gormConfig.stats {
		ag.gormMiddleware(ctx, request)
	}
}

// MiddlewareXml 中间件
func (ag *ApiGorm) MiddlewareXml(ctx context.Context, request gorequest.Response) {
	if ag.gormConfig.stats {
		ag.gormMiddlewareXml(ctx, request)
	}
}

// MiddlewareCustom 中间件
func (ag *ApiGorm) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	if ag.gormConfig.stats {
		ag.gormMiddlewareCustom(ctx, api, request)
	}
}
