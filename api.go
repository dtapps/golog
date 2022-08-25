package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"os"
	"runtime"
)

type ApiClientConfig struct {
	GormClient *dorm.GormClient // 数据库驱动
	TableName  string           // 表名
	LogClient  *ZapLog          // 日志驱动
	LogDebug   bool             // 日志开关
}

// ApiClient 接口
type ApiClient struct {
	gormClient  *dorm.GormClient  // 数据库驱动
	mongoClient *dorm.MongoClient // 数据库驱动
	logClient   *ZapLog           // 日志驱动
	config      struct {
		tableName string // 表名
		insideIp  string // 内网ip
		hostname  string // 主机名
		goVersion string // go版本
		logDebug  bool   // 日志开关
	}
	mongoConfig struct {
		databaseName   string // 库名
		collectionName string // 表名
		insideIp       string // 内网ip
		hostname       string // 主机名
		goVersion      string // go版本
		debug          bool   // 日志开关
	}
}

// NewApiClient 创建接口实例化
func NewApiClient(config *ApiClientConfig) (*ApiClient, error) {

	c := &ApiClient{}

	c.gormClient = config.GormClient
	c.config.tableName = config.TableName

	c.logClient = config.LogClient
	c.config.logDebug = config.LogDebug

	if c.gormClient.Db == nil {
		return nil, errors.New("驱动不能为空")
	}

	if c.config.tableName == "" {
		return nil, errors.New("表名不能为空")
	}

	err := c.gormClient.Db.Table(c.config.tableName).AutoMigrate(&apiPostgresqlLog{})
	if err != nil {
		return nil, errors.New("创建表失败：" + err.Error())
	}

	hostname, _ := os.Hostname()

	c.config.hostname = hostname
	c.config.insideIp = goip.GetInsideIp(context.Background())
	c.config.goVersion = runtime.Version()

	return c, nil
}

// NewApiMongoClient 创建接口实例化
func NewApiMongoClient(mongoClient *dorm.MongoClient, databaseName string, collectionName string, debug bool) (*ApiClient, error) {

	c := &ApiClient{}

	if mongoClient.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	c.mongoClient = mongoClient

	if databaseName == "" {
		return nil, errors.New("没有设置库名")
	}
	c.mongoConfig.databaseName = databaseName

	if collectionName == "" {
		return nil, errors.New("没有设置表名")
	}
	c.mongoConfig.collectionName = collectionName

	c.mongoConfig.debug = debug

	hostname, _ := os.Hostname()

	c.mongoConfig.hostname = hostname
	c.mongoConfig.insideIp = goip.GetInsideIp(context.Background())
	c.mongoConfig.goVersion = runtime.Version()

	return c, nil
}
