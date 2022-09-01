package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
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
	log struct {
		gorm  bool // 日志开关
		mongo bool // 日志开关
	}
}

// client 数据库服务
// string 表名
type apiGormClientFun func() (*dorm.GormClient, string)

// client 数据库服务
// string 库名
// string 表名
type apiMongoClientFun func() (*dorm.MongoClient, string, string)

// ApiClientConfig 接口实例配置
type ApiClientConfig struct {
	GormClientFun  apiGormClientFun  // 日志配置
	MongoClientFun apiMongoClientFun // 日志配置
	Debug          bool              // 日志开关
}

// NewApiClient 创建接口实例化
// client 数据库服务
// tableName 表名
func NewApiClient(config *ApiClientConfig) (*ApiClient, error) {

	var ctx = context.Background()

	c := &ApiClient{}

	gormClient, gormTableName := config.GormClientFun()
	mongoClient, mongoDatabaseName, mongoCollectionName := config.MongoClientFun()

	if (gormClient == nil || gormClient.Db == nil) || (mongoClient == nil || mongoClient.Db == nil) {
		return nil, errors.New("没有设置驱动")
	}

	hostname, _ := os.Hostname()

	if gormClient != nil || gormClient.Db != nil {

		c.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		}
		c.gormConfig.tableName = gormTableName

		c.gormConfig.debug = config.Debug

		err := c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&apiPostgresqlLog{})
		if err != nil {
			return nil, errors.New("创建表失败：" + err.Error())
		}

		c.gormConfig.hostname = hostname
		c.gormConfig.insideIp = goip.GetInsideIp(ctx)
		c.gormConfig.goVersion = runtime.Version()

		c.log.gorm = true

	}

	if mongoClient != nil || mongoClient.Db != nil {

		c.mongoClient = mongoClient

		if mongoDatabaseName == "" {
			return nil, errors.New("没有设置库名")
		}
		c.mongoConfig.databaseName = mongoDatabaseName

		if mongoCollectionName == "" {
			return nil, errors.New("没有设置表名")
		}
		c.mongoConfig.collectionName = mongoCollectionName

		c.mongoConfig.debug = config.Debug

		c.mongoConfig.hostname = hostname
		c.mongoConfig.insideIp = goip.GetInsideIp(ctx)
		c.mongoConfig.goVersion = runtime.Version()

		c.log.mongo = true

		// 创建时间序列集合
		c.mongoCreateCollection()

		// 创建索引
		c.mongoCreateIndexes()

	}

	return c, nil
}

// Middleware 中间件
func (c *ApiClient) Middleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	if c.log.gorm {
		c.GormMiddleware(ctx, request, sdkVersion)
	}
	if c.log.mongo {
		c.MongoMiddleware(ctx, request, sdkVersion)
	}
}

// MiddlewareXml 中间件
func (c *ApiClient) MiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	if c.log.gorm {
		c.GormMiddlewareXml(ctx, request, sdkVersion)
	}
	if c.log.mongo {
		c.MongoMiddlewareXml(ctx, request, sdkVersion)
	}
}

// MiddlewareCustom 中间件
func (c *ApiClient) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	if c.log.gorm {
		c.GormMiddlewareCustom(ctx, api, request, sdkVersion)
	}
	if c.log.mongo {
		c.MongoMiddlewareCustom(ctx, api, request, sdkVersion)
	}
}
