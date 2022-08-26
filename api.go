package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"os"
	"runtime"
)

// ApiClient 接口
type ApiClient struct {
	gormClient  *dorm.GormClient  // 数据库驱动
	mongoClient *dorm.MongoClient // 数据库驱动
	gormConfig  struct {
		tableName string // 表名
		insideIp  string // 内网ip
		hostname  string // 主机名
		goVersion string // go版本
		debug     bool   // 日志开关
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

// NewApiGormClient 创建接口实例化
// client 数据库服务
// tableName 表名
func NewApiGormClient(gormClientFun func() (client *dorm.GormClient, tableName string), debug bool) (*ApiClient, error) {

	c := &ApiClient{}

	client, tableName := gormClientFun()

	if client == nil || client.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	c.gormClient = client

	if tableName == "" {
		return nil, errors.New("没有设置表名")
	}
	c.gormConfig.tableName = tableName

	c.gormConfig.debug = debug

	err := c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&apiPostgresqlLog{})
	if err != nil {
		return nil, errors.New("创建表失败：" + err.Error())
	}

	hostname, _ := os.Hostname()

	c.gormConfig.hostname = hostname
	c.gormConfig.insideIp = goip.GetInsideIp(context.Background())
	c.gormConfig.goVersion = runtime.Version()

	return c, nil
}

// NewApiMongoClient 创建接口实例化
// client 数据库服务
// databaseName 库名
// collectionName 表名
func NewApiMongoClient(mongoClientFun func() (client *dorm.MongoClient, databaseName string, collectionName string), debug bool) (*ApiClient, error) {

	c := &ApiClient{}

	client, databaseName, collectionName := mongoClientFun()

	if client == nil || client.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	c.mongoClient = client

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
