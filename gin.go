package golog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/goxml"
	"io/ioutil"
	"net"
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
	log struct {
		gorm  bool // 日志开关
		mongo bool // 日志开关
	}
}

// client 数据库服务
// string 表名
type ginGormClientFun func() (*dorm.GormClient, string)

// client 数据库服务
// string 库名
// string 表名
type ginMongoClientFun func() (*dorm.MongoClient, string, string)

// GinClientConfig 框架实例配置
type GinClientConfig struct {
	IpService      *goip.Client      // ip服务
	GormClientFun  ginGormClientFun  // 日志配置
	MongoClientFun apiMongoClientFun // 日志配置
	Debug          bool              // 日志开关
}

// NewGinClient 创建框架实例化
// client 数据库服务
// tableName 表名
// ipService ip服务
func NewGinClient(config *GinClientConfig) (*GinClient, error) {

	var ctx = context.Background()

	c := &GinClient{}

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

		c.ipService = config.IpService

		err := c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&ginPostgresqlLog{})
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

		c.ipService = config.IpService

		c.mongoConfig.hostname = hostname
		c.mongoConfig.insideIp = goip.GetInsideIp(ctx)
		c.mongoConfig.goVersion = runtime.Version()

		c.log.mongo = true
	}

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

// Middleware 中间件
func (c *GinClient) Middleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		// 开始时间
		startTime := gotime.Current().TimestampWithMillisecond()
		requestTime := gotime.Current().Time

		// 获取
		data, _ := ioutil.ReadAll(ginCtx.Request.Body)

		// 复用
		ginCtx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
		ginCtx.Writer = blw

		// 处理请求
		ginCtx.Next()

		// 响应
		responseCode := ginCtx.Writer.Status()
		responseBody := blw.body.String()

		//结束时间
		endTime := gotime.Current().TimestampWithMillisecond()

		go func() {

			var dataJson = true

			// 解析请求内容
			var xmlBody map[string]string
			var jsonBody map[string]interface{}

			// 判断是否有内容
			if len(data) > 0 {
				err := json.Unmarshal(data, &jsonBody)
				if len(jsonBody) <= 0 {
					dataJson = false
					xmlBody = goxml.XmlDecode(string(data))
				}

				if err != nil {
					dataJson = false
					xmlBody = goxml.XmlDecode(string(data))
				}
			}

			clientIp := gorequest.ClientIp(ginCtx.Request)

			requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp := "", "", "", "", ""
			if c.ipService != nil {
				if net.ParseIP(clientIp).To4() != nil {
					// IPv4
					_, info := c.ipService.Ipv4(clientIp)
					requestClientIpCountry = info.Country
					requestClientIpRegion = info.Region
					requestClientIpProvince = info.Province
					requestClientIpCity = info.City
					requestClientIpIsp = info.ISP
				} else if net.ParseIP(clientIp).To16() != nil {
					// IPv6
					info := c.ipService.Ipv6(clientIp)
					requestClientIpCountry = info.Country
					requestClientIpProvince = info.Province
					requestClientIpCity = info.City
				}
			}

			var traceId = gotrace_id.GetGinTraceId(ginCtx)

			// 记录
			if c.log.gorm {
				if dataJson {
					c.gormRecordJson(ginCtx, traceId, requestTime, jsonBody, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				} else {
					c.gormRecordXml(ginCtx, traceId, requestTime, xmlBody, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				}
			}
			// 记录
			if c.log.mongo {
				if dataJson {
					c.mongoRecordJson(ginCtx, traceId, requestTime, jsonBody, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				} else {
					c.mongoRecordXml(ginCtx, traceId, requestTime, xmlBody, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				}
			}
		}()
	}
}
