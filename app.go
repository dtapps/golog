package golog

import (
	"go.dtapp.net/goip"
	"gorm.io/gorm"
	"os"
	"runtime"
	"strconv"
	"strings"
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
	a.Api.hostname, _ = os.Hostname()
	a.Api.insideIp = goip.GetInsideIp()
	goVersion, _ := strconv.ParseFloat(strings.TrimPrefix(runtime.Version(), "go"), 64)
	a.Api.goVersion = goVersion
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
	a.Gin.hostname, _ = os.Hostname()
	a.Gin.insideIp = goip.GetInsideIp()
	goVersion, _ := strconv.ParseFloat(strings.TrimPrefix(runtime.Version(), "go"), 64)
	a.Gin.goVersion = goVersion
	a.Gin.AutoMigrate()
}
