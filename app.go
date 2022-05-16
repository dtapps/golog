package golog

import (
	"go.dtapp.net/goip"
	"gorm.io/gorm"
)

type App struct {
	Gin       gin      // 框架日志
	Api       api      // 接口日志
	Pgsql     *gorm.DB // pgsql数据库
	TableName string   // 日志表名
}

// InitClientApi 接口实例化
func (a *App) InitClientApi() {
	if a.Pgsql == nil {
		panic("驱动不正常")
	}
	if a.TableName == "" {
		panic("表名不能为空")
	}
	a.Api.db = a.Pgsql
	a.Api.tableName = a.TableName
	a.Api.outsideIp = goip.GetOutsideIp()
	a.Api.insideIp = goip.GetInsideIp()
	a.Api.AutoMigrate()
}

// InitClientGin 框架实例化
func (a *App) InitClientGin() {
	if a.Pgsql == nil {
		panic("驱动不正常")
	}
	if a.TableName == "" {
		panic("表名不能为空")
	}
	a.Gin.db = a.Pgsql
	a.Gin.tableName = a.TableName
	a.Gin.outsideIp = goip.GetOutsideIp()
	a.Gin.insideIp = goip.GetInsideIp()
	a.Gin.AutoMigrate()
}
