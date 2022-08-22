package golog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"os"
	"runtime"
)

type GinClientConfig struct {
	GormClient *dorm.GormClient // 数据库驱动
	IpService  *goip.Client     // ip服务
	TableName  string           // 表名
	LogClient  *ZapLog          // 日志驱动
	LogDebug   bool             // 日志开关
}

// GinClient 框架
type GinClient struct {
	gormClient  *dorm.GormClient  // 数据库驱动
	mongoClient *dorm.MongoClient // 数据库驱动
	ipService   *goip.Client      // ip服务
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

// NewGinClient 创建框架实例化
func NewGinClient(config *GinClientConfig) (*GinClient, error) {

	c := &GinClient{}

	c.gormClient = config.GormClient
	c.config.tableName = config.TableName

	c.logClient = config.LogClient
	c.config.logDebug = config.LogDebug

	c.ipService = config.IpService

	if c.gormClient.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	if c.config.tableName == "" {
		return nil, errors.New("没有设置表名")
	}

	err := c.gormClient.Db.Table(c.config.tableName).AutoMigrate(&ginPostgresqlLog{})
	if err != nil {
		return nil, errors.New("创建表失败：" + err.Error())
	}

	hostname, _ := os.Hostname()

	c.config.hostname = hostname
	c.config.insideIp = goip.GetInsideIp(context.Background())
	c.config.goVersion = runtime.Version()

	return c, nil
}

// NewGinMongoClient 创建框架实例化
// mongoClient 数据库服务
// ipService ip服务
// databaseName 库名
// collectionName 表名
func NewGinMongoClient(mongoClient *dorm.MongoClient, ipService *goip.Client, databaseName string, collectionName string, debug bool) (*GinClient, error) {

	c := &GinClient{}

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

	c.ipService = ipService

	c.mongoConfig.debug = debug

	hostname, _ := os.Hostname()

	c.mongoConfig.hostname = hostname
	c.mongoConfig.insideIp = goip.GetInsideIp(context.Background())
	c.mongoConfig.goVersion = runtime.Version()

	return c, nil
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (c *GinClient) jsonUnmarshal(data string) (result interface{}) {
	_ = json.Unmarshal([]byte(data), &result)
	return
}
