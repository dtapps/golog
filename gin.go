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

// GinClient 框架
type GinClient struct {
	gormClient  *dorm.GormClient  // 数据库驱动
	mongoClient *dorm.MongoClient // 数据库驱动
	ipService   *goip.Client      // ip服务
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

// NewGinGormClient 创建框架实例化
// client 数据库服务
// tableName=表名
// ipService ip服务
func NewGinGormClient(gormClientFun func() (client *dorm.GormClient, tableName string), ipService *goip.Client, debug bool) (*GinClient, error) {

	c := &GinClient{}

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

	c.ipService = ipService

	err := c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&ginPostgresqlLog{})
	if err != nil {
		return nil, errors.New("创建表失败：" + err.Error())
	}

	hostname, _ := os.Hostname()

	c.gormConfig.hostname = hostname
	c.gormConfig.insideIp = goip.GetInsideIp(context.Background())
	c.gormConfig.goVersion = runtime.Version()

	return c, nil
}

// NewGinMongoClient 创建框架实例化
// client 数据库服务
// databaseName 库名
// collectionName 表名
// ipService ip服务
func NewGinMongoClient(mongoClientFun func() (client *dorm.MongoClient, databaseName string, collectionName string), ipService *goip.Client, debug bool) (*GinClient, error) {

	c := &GinClient{}

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
