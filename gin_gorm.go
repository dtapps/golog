package golog

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"gorm.io/gorm"
	"io/ioutil"
)

// GinGorm 框架日志
type GinGorm struct {
	gormClient *gorm.DB // 数据库驱动
	config     struct {
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
		systemInsideIp      string  // 内网ip
		systemOutsideIp     string  // 外网ip
		goVersion           string  // go版本
		sdkVersion          string  // sdk版本
	}
	gormConfig struct {
		stats     bool   // 状态
		tableName string // 表名
	}
}

// GinGormFun *GinGorm 框架日志驱动
type GinGormFun func() *GinGorm

// NewGinGorm 创建框架实例化
func NewGinGorm(ctx context.Context, systemOutsideIp string, gormClient *gorm.DB, gormTableName string) (*GinGorm, error) {

	gg := &GinGorm{}

	// 配置信息
	if systemOutsideIp == "" {
		return nil, errors.New("没有设置外网IP")
	}
	gg.setConfig(ctx, systemOutsideIp)

	if gormClient == nil {
		gg.gormConfig.stats = false
	} else {

		gg.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			gg.gormConfig.tableName = gormTableName
		}

		// 创建模型
		gg.gormAutoMigrate(ctx)

		gg.gormConfig.stats = true
	}

	return gg, nil
}

type bodyGormWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyGormWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyGormWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (gg *GinGorm) jsonUnmarshal(data string) (result interface{}) {
	_ = gojson.Unmarshal([]byte(data), &result)
	return
}

// Middleware 中间件
func (gg *GinGorm) Middleware() gin.HandlerFunc {
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
		var dataMap map[string]interface{}
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

		blw := &bodyGormWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
		ginCtx.Writer = blw

		// 处理请求
		ginCtx.Next()

		// 响应
		responseCode := ginCtx.Writer.Status()
		responseBody := blw.body.String()

		// 结束时间
		endTime := gotime.Current().TimestampWithMillisecond()
		responseTime := gotime.Current().Time

		go func() {

			// 记录
			gg.recordJson(ginCtx, requestTime, requestBody, responseTime, responseCode, responseBody, endTime-startTime, gorequest.ClientIp(ginCtx.Request))

		}()
	}
}
