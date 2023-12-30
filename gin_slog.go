package golog

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"io/ioutil"
)

// GinSLog 框架日志
type GinSLog struct {
	slog struct {
		status bool  // 状态
		client *SLog // 日志服务
	}
}

// GinSLogFun  框架日志驱动
type GinSLogFun func() *GinSLog

// NewGinSLog 创建框架实例化
func NewGinSLog(ctx context.Context) (*GinSLog, error) {
	c := &GinSLog{}
	return c, nil
}

type bodySLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodySLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodySLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (gl *GinSLog) jsonUnmarshal(data string) (result interface{}) {
	_ = gojson.Unmarshal([]byte(data), &result)
	return
}

// Middleware 中间件
func (gl *GinSLog) Middleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		// 开始时间
		startTime := gotime.Current().TimestampWithMillisecond()
		requestTime := gotime.Current().Time

		// 获取全部内容
		paramsBody := gorequest.NewParams()
		queryParams := ginCtx.Request.URL.Query() // 请求URL参数
		for key, values := range queryParams {
			for _, value := range values {
				paramsBody.Set(key, value)
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
			paramsBody.Set(key, value)
		}

		// 重新赋值
		ginCtx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))

		blw := &bodySLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
		ginCtx.Writer = blw

		// 处理请求
		ginCtx.Next()

		// 响应
		responseCode := ginCtx.Writer.Status()
		responseBody := blw.body.String()

		// 结束时间
		endTime := gotime.Current().TimestampWithMillisecond()

		go func() {

			requestIp := gorequest.ClientIp(ginCtx.Request)

			var traceId = gotrace_id.GetGinTraceId(ginCtx)

			// 记录
			gl.recordJson(ginCtx, traceId, requestTime, paramsBody, responseCode, responseBody, startTime, endTime, requestIp)

		}()
	}
}
