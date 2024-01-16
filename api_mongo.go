package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gorequest"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApiMongo 接口日志
type ApiMongo struct {
	mongoClient *mongo.Client // 数据库驱动
	config      struct {
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
	mongoConfig struct {
		stats          bool   // 状态
		databaseName   string // 库名
		collectionName string // 集合名
	}
	slog struct {
		status bool  // 状态
		client *SLog // 日志服务
	}
}

// ApiMongoFun 接口日志驱动
type ApiMongoFun func() *ApiMongo

// NewApiMongo 创建接口实例化
func NewApiMongo(ctx context.Context, systemOutsideIp string, mongoClient *mongo.Client, mongoDatabaseName string, mongoCollectionName string) (*ApiMongo, error) {

	am := &ApiMongo{}

	// 配置信息
	if systemOutsideIp == "" {
		return nil, errors.New("没有设置外网IP")
	}
	am.setConfig(ctx, systemOutsideIp)

	if mongoClient == nil {
		am.mongoConfig.stats = false
	} else {

		am.mongoClient = mongoClient

		if mongoDatabaseName == "" {
			return nil, errors.New("没有设置库名")
		} else {
			am.mongoConfig.databaseName = mongoDatabaseName
		}

		if mongoCollectionName == "" {
			return nil, errors.New("没有设置集合名")
		} else {
			am.mongoConfig.collectionName = mongoCollectionName
		}

		am.mongoConfig.stats = true

		// 创建时间序列集合
		am.mongoCreateCollection(ctx)

		// 创建索引
		am.mongoCreateIndexes(ctx)

	}

	return am, nil
}

// Middleware 中间件
func (am *ApiMongo) Middleware(ctx context.Context, request gorequest.Response) {
	if am.mongoConfig.stats {
		am.mongoMiddleware(ctx, request)
	}
}

// MiddlewareXml 中间件
func (am *ApiMongo) MiddlewareXml(ctx context.Context, request gorequest.Response) {
	if am.mongoConfig.stats {
		am.mongoMiddlewareXml(ctx, request)
	}
}

// MiddlewareCustom 中间件
func (am *ApiMongo) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	if am.mongoConfig.stats {
		am.mongoMiddlewareCustom(ctx, api, request)
	}
}
