package golog

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// GinGorm 框架日志
type GinGorm struct {
	gormClient *gorm.DB // 数据库驱动
	config     struct {
		GoVersion  string // go版本
		SdkVersion string // sdk版本
		system     struct {
			SystemVersion  string  `json:"system_version"`   // 系统版本
			SystemOs       string  `json:"system_os"`        // 系统类型
			SystemArch     string  `json:"system_arch"`      // 系统内核
			SystemInsideIP string  `json:"system_inside_ip"` // 内网IP
			SystemCpuModel string  `json:"system_cpu_model"` // CPU型号
			SystemCpuCores int     `json:"system_cpu_cores"` // CPU核数
			SystemCpuMhz   float64 `json:"system_cpu_mhz"`   // CPU兆赫
		}
	}
	gormConfig struct {
		stats     bool   // 状态
		tableName string // 表名
	}
}

// GinGormFun *GinGorm 框架日志驱动
type GinGormFun func() *GinGorm

// NewGinGorm 创建框架实例化
func NewGinGorm(ctx context.Context, gormClient *gorm.DB, gormTableName string) (*GinGorm, error) {

	gg := &GinGorm{}
	gg.setConfig(ctx)

	if gormClient == nil {
		gg.gormConfig.stats = false
	} else {

		gg.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			gg.gormConfig.tableName = gormTableName
		}

		gg.gormConfig.stats = true

		// 创建模型
		gg.gormAutoMigrate(ctx)

	}

	return gg, nil
}

// 定义一个自定义的 ResponseWriter
type ginGormBodyWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	status  int
	headers http.Header
}

// 实现 http.ResponseWriter 的 Write 方法
func (w ginGormBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString 实现 http.ResponseWriter 的 WriteString 方法
//func (w ginGormBodyWriter) WriteString(s string) (int, error) {
//	w.body.WriteString(s)
//	return w.ResponseWriter.WriteString(s)
//}

// WriteHeader 实现 http.ResponseWriter 的 WriteHeader 方法
func (w ginGormBodyWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Header 实现 http.ResponseWriter 的 Header 方法
func (w ginGormBodyWriter) Header() http.Header {
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	return w.headers
}

func (gg *GinGorm) jsonUnmarshal(data string) (result any) {
	_ = gojson.Unmarshal([]byte(data), &result)
	return
}

// Middleware 中间件
func (gg *GinGorm) Middleware() gin.HandlerFunc {
	return func(g *gin.Context) {

		// 开始时间
		start := time.Now().UTC()

		// 模型
		var log = ginGormLog{}

		// 请求时间
		log.RequestTime = gotime.Current().Time

		// 创建自定义的 ResponseWriter 并替换原有的
		blw := &ginGormBodyWriter{
			ResponseWriter: g.Writer,
			body:           bytes.NewBufferString(""),
		}
		g.Writer = blw

		// 处理请求
		g.Next()

		// 结束时间
		end := time.Now().UTC()

		// 请求消耗时长
		log.RequestCostTime = end.Sub(start).Milliseconds()

		// 响应时间
		log.ResponseTime = gotime.Current().Time

		// 请求编号
		log.RequestID = gotrace_id.GetGinTraceId(g)

		// 请求主机
		log.RequestHost = g.Request.Host

		// 请求地址
		log.RequestPath = gourl.UriFilterExcludeQueryString(g.Request.RequestURI)

		// 请求参数
		log.RequestQuery = gojson.JsonEncodeNoError(gojson.ParseQueryString(string(g.Request.RequestURI)))

		// 请求方式
		log.RequestMethod = g.Request.Method

		// 请求协议
		log.RequestScheme = g.Request.Proto

		// 请求类型
		log.RequestContentType = g.ContentType()

		// 请求IP
		log.RequestClientIP = g.ClientIP()

		// 请求UA
		log.RequestUserAgent = g.Request.UserAgent()

		// 请求头
		log.RequestHeader = gojson.JsonEncodeNoError(g.Request.Header)

		// 响应头
		log.ResponseHeader = gojson.JsonEncodeNoError(blw.Header())

		// 响应状态
		log.ResponseStatusCode = g.Writer.Status()

		// 响应内容
		if gojson.IsValidJSON(blw.body.String()) {
			log.ResponseBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(blw.body.String()))
		} else {
			log.ResponseBody = blw.body.String()
		}

		go func() {
			gg.gormRecord(g, log)
		}()
	}
}
