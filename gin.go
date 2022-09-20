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
	"io/ioutil"
)

// GinClientFun *GinClient 驱动
type GinClientFun func() *GinClient

// GinClient 框架
type GinClient struct {
	gormClient  *dorm.GormClient  // 数据库驱动
	mongoClient *dorm.MongoClient // 数据库驱动
	ipService   *goip.Client      // ip服务
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

// GinClientConfig 框架实例配置
type GinClientConfig struct {
	IpService      *goip.Client                  // ip服务
	GormClientFun  dorm.GormClientTableFun       // 日志配置
	MongoClientFun dorm.MongoClientCollectionFun // 日志配置
	Debug          bool                          // 日志开关
	ZapLog         *ZapLog                       // 日志服务
}

// NewGinClient 创建框架实例化
func NewGinClient(config *GinClientConfig) (*GinClient, error) {

	var ctx = context.Background()

	c := &GinClient{}

	c.zapLog = config.ZapLog

	c.logDebug = config.Debug

	c.ipService = config.IpService

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
			var jsonBody map[string]interface{}

			// 判断是否有内容
			if len(data) > 0 {
				err := json.Unmarshal(data, &jsonBody)
				if err != nil {
					dataJson = false
				}
			}

			clientIp := gorequest.ClientIp(ginCtx.Request)

			var requestClientIpCountry string
			var requestClientIpProvince string
			var requestClientIpCity string
			var requestClientIpIsp string
			var requestClientIpLocationLatitude float64
			var requestClientIpLocationLongitude float64
			if c.ipService != nil {
				info := c.ipService.Analyse(clientIp)
				requestClientIpCountry = info.Country
				requestClientIpProvince = info.Province
				requestClientIpCity = info.City
				requestClientIpIsp = info.Isp
				requestClientIpLocationLatitude = info.LocationLatitude
				requestClientIpLocationLongitude = info.LocationLongitude
			}

			var traceId = gotrace_id.GetGinTraceId(ginCtx)

			// 记录
			if c.gormConfig.stats {
				if dataJson {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.Middleware]准备使用{gormRecordJson}保存数据：%s", data)
					}
					c.gormRecordJson(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp, requestClientIpLocationLatitude, requestClientIpLocationLongitude)
				} else {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.Middleware]准备使用{gormRecordXml}保存数据：%s", data)
					}
					c.gormRecordXml(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp, requestClientIpLocationLatitude, requestClientIpLocationLongitude)
				}
			}
			if c.mongoConfig.stats {
				if dataJson {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.Middleware]准备使用{mongoRecordJson}保存数据：%s", data)
					}
					c.mongoRecordJson(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp, requestClientIpLocationLatitude, requestClientIpLocationLongitude)
				} else {
					if c.logDebug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.Middleware]准备使用{mongoRecordXml}保存数据：%s", data)
					}
					c.mongoRecordXml(ginCtx, traceId, requestTime, data, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp, requestClientIpLocationLatitude, requestClientIpLocationLongitude)
				}
			}
		}()
	}
}
