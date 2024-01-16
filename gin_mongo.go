package golog

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
)

// GinMongo 框架日志
type GinMongo struct {
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

// GinMongoFun *GinMongo 框架日志驱动
type GinMongoFun func() *GinMongo

// NewGinMongo 创建框架实例化
func NewGinMongo(ctx context.Context, systemOutsideIp string, mongoClient *mongo.Client, mongoDatabaseName string, mongoCollectionName string) (*GinMongo, error) {

	gm := &GinMongo{}

	// 配置信息
	if systemOutsideIp == "" {
		return nil, errors.New("没有设置外网IP")
	}
	gm.setConfig(ctx, systemOutsideIp)

	if mongoClient == nil {
		gm.mongoConfig.stats = false
	} else {

		gm.mongoClient = mongoClient

		if mongoDatabaseName == "" {
			return nil, errors.New("没有设置库名")
		} else {
			gm.mongoConfig.databaseName = mongoDatabaseName
		}

		if mongoCollectionName == "" {
			return nil, errors.New("没有设置集合名")
		} else {
			gm.mongoConfig.collectionName = mongoCollectionName
		}

		gm.mongoConfig.stats = true

		// 创建时间序列集合
		gm.mongoCreateCollection(ctx)

		// 创建索引
		gm.mongoCreateIndexes(ctx)

	}

	return gm, nil
}

type ginMongoBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w ginMongoBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ginMongoBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (gm *GinMongo) jsonUnmarshal(data string) (result any) {
	_ = gojson.Unmarshal([]byte(data), &result)
	return
}

// Middleware 中间件
func (gm *GinMongo) Middleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		// 开始时间
		startTime := gotime.Current().TimestampWithMillisecond()
		requestTime := gotime.Current().Time

		// 获取全部内容
		requestBody := gorequest.NewParams()
		queryParams := ginCtx.Request.URL.Query() // 请求URL参数
		for key, values := range queryParams {
			for _, value := range values {
				requestBody.Set(key, value)
			}
		}
		var dataMap map[string]any
		rawData, _ := ginCtx.GetRawData() // 请求内容参数
		if gojson.IsValidJSON(string(rawData)) {
			dataMap = gojson.JsonDecodeNoError(string(rawData))
		} else {
			dataMap = gojson.ParseQueryString(string(rawData))
		}
		for key, value := range dataMap {
			requestBody.Set(key, value)
		}

		// 重新赋值
		ginCtx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))

		blw := &ginMongoBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
		ginCtx.Writer = blw

		// 处理请求
		ginCtx.Next()

		// 响应
		responseCode := ginCtx.Writer.Status()
		responseBody := gojson.JsonDecodeNoError(blw.body.String())

		// 结束时间
		endTime := gotime.Current().TimestampWithMillisecond()
		responseTime := gotime.Current().Time

		go func() {

			// 记录
			gm.recordJson(ginCtx, requestTime, requestBody, responseTime, responseCode, responseBody, endTime-startTime, gorequest.ClientIp(ginCtx.Request))

		}()
	}
}
