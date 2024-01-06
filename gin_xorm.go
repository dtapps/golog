package golog

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"io/ioutil"
	"xorm.io/xorm"
)

// GinXorm 框架日志
type GinXorm struct {
	xormClient *xorm.Engine // 数据库驱动
	config     struct {
		systemHostname      string  // 主机名
		systemOs            string  // 系统类型
		systemVersion       string  // 系统版本
		systemKernel        string  // 系统内核
		systemKernelVersion string  // 系统内核版本
		systemBootTime      uint64  // 系统开机时间
		cpuCores            int     // CPU核数
		cpuModelName        string  // CPU型号名称
		cpuMhz              float64 // CPU兆赫
		systemInsideIp      string  // 内网ip
		systemOutsideIp     string  // 外网ip
		goVersion           string  // go版本
		sdkVersion          string  // sdk版本
	}
	xormConfig struct {
		stats     bool   // 状态
		tableName string // 表名
	}
}

// GinXormFun *GinXorm 框架日志驱动
type GinXormFun func() *GinXorm

// NewGinXorm 创建框架实例化
func NewGinXorm(ctx context.Context, systemOutsideIp string, xormClient *xorm.Engine, xormTableName string) (*GinXorm, error) {

	gg := &GinXorm{}

	// 配置信息
	if systemOutsideIp == "" {
		return nil, errors.New("没有设置外网IP")
	}
	gg.setConfig(ctx, systemOutsideIp)

	if xormClient == nil {
		gg.xormConfig.stats = false
	} else {

		gg.xormClient = xormClient

		if xormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			gg.xormConfig.tableName = xormTableName
		}

		// 创建模型
		gg.xormSync(ctx)

		gg.xormConfig.stats = true
	}

	return gg, nil
}

type bodyXormWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyXormWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyXormWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (gg *GinXorm) jsonUnmarshal(data string) (result interface{}) {
	_ = gojson.Unmarshal([]byte(data), &result)
	return
}

// Middleware 中间件
func (gg *GinXorm) Middleware() gin.HandlerFunc {
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

		blw := &bodyXormWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
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
