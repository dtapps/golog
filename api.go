package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
)

// ApiClientFun *ApiClient 驱动
type ApiClientFun func() *ApiClient

// ApiClient 接口
type ApiClient struct {
	gormClient  *dorm.GormClient  // 数据库驱动
	mongoClient *dorm.MongoClient // 数据库驱动
	zapLog      *ZapLog           // 日志服务
	logDebug    bool              // 日志开关
	config      struct {
		systemHostName  string // 主机名
		systemInsideIp  string // 内网ip
		systemOs        string // 系统类型
		systemArch      string // 系统架构
		goVersion       string // go版本
		sdkVersion      string // sdk版本
		systemOutsideIp string // 外网ip
	}
	gormConfig struct {
		stats     bool   // 状态
		tableName string // 表名
	}
	mongoConfig struct {
		stats          bool   // 状态
		databaseName   string // 库名
		collectionName string // 表名
	}
}

// ApiClientConfig 接口实例配置
type ApiClientConfig struct {
	GormClientFun  dorm.GormClientTableFun       // 日志配置
	MongoClientFun dorm.MongoClientCollectionFun // 日志配置
	Debug          bool                          // 日志开关
	ZapLog         *ZapLog                       // 日志服务
	CurrentIp      string                        // 当前ip
}

// NewApiClient 创建接口实例化
func NewApiClient(config *ApiClientConfig) (*ApiClient, error) {

	var ctx = context.Background()

	c := &ApiClient{}

	c.zapLog = config.ZapLog

	c.logDebug = config.Debug

	if config.CurrentIp == "" {
		config.CurrentIp = goip.GetOutsideIp(ctx)
	}
	if config.CurrentIp != "" && config.CurrentIp != "0.0.0.0" {
		c.config.systemOutsideIp = config.CurrentIp
	}

	if c.config.systemOutsideIp == "" {
		return nil, currentIpNoConfig
	}

	// 配置信息
	c.setConfig(ctx)

	gormClient, gormTableName := config.GormClientFun()
	mongoClient, mongoDatabaseName, mongoCollectionName := config.MongoClientFun()

	if (gormClient == nil || gormClient.Db == nil) || (mongoClient == nil || mongoClient.Db == nil) {
		return nil, dbClientFunNoConfig
	}

	// 配置关系数据库
	if gormClient != nil || gormClient.Db != nil {

		c.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			c.gormConfig.tableName = gormTableName
		}

		// 创建模型
		c.gormAutoMigrate(ctx)

		c.gormConfig.stats = true
	}

	// 配置非关系数据库
	if mongoClient != nil || mongoClient.Db != nil {

		c.mongoClient = mongoClient

		if mongoDatabaseName == "" {
			return nil, errors.New("没有设置库名")
		} else {
			c.mongoConfig.databaseName = mongoDatabaseName
		}

		if mongoCollectionName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			c.mongoConfig.collectionName = mongoCollectionName
		}

		// 创建时间序列集合
		c.mongoCreateCollection(ctx)

		// 创建索引
		c.mongoCreateIndexes(ctx)

		c.mongoConfig.stats = true
	}

	return c, nil
}

// Middleware 中间件
func (c *ApiClient) Middleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	if c.gormConfig.stats {
		c.gormMiddleware(ctx, request, sdkVersion)
	}
	if c.mongoConfig.stats {
		c.mongoMiddleware(ctx, request, sdkVersion)
	}
}

// MiddlewareXml 中间件
func (c *ApiClient) MiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	if c.gormConfig.stats {
		c.gormMiddlewareXml(ctx, request, sdkVersion)
	}
	if c.mongoConfig.stats {
		c.mongoMiddlewareXml(ctx, request, sdkVersion)
	}
}

// MiddlewareCustom 中间件
func (c *ApiClient) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	if c.gormConfig.stats {
		c.gormMiddlewareCustom(ctx, api, request, sdkVersion)
	}
	if c.mongoConfig.stats {
		c.mongoMiddlewareCustom(ctx, api, request, sdkVersion)
	}
}
