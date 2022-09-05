package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"os"
	"runtime"
	"time"
	"unicode/utf8"
)

// ApiGormClientConfig 接口实例配置
type ApiGormClientConfig struct {
	GormClientFun apiGormClientFun // 日志配置
	Debug         bool             // 日志开关
	ZapLog        *ZapLog          // 日志服务
}

// NewApiGormClient 创建接口实例化
// client 数据库服务
// tableName 表名
func NewApiGormClient(config *ApiGormClientConfig) (*ApiClient, error) {

	var ctx = context.Background()

	c := &ApiClient{}

	c.zapLog = config.ZapLog

	client, tableName := config.GormClientFun()

	if client == nil || client.Db == nil {
		return nil, errors.New("没有设置驱动")
	}

	c.gormClient = client

	if tableName == "" {
		return nil, errors.New("没有设置表名")
	}
	c.gormConfig.tableName = tableName

	c.gormConfig.debug = config.Debug

	err := c.gormAutoMigrate()
	if err != nil {
		return nil, errors.New("创建表失败：" + err.Error())
	}

	hostname, _ := os.Hostname()

	c.gormConfig.hostname = hostname
	c.gormConfig.insideIp = goip.GetInsideIp(ctx)
	c.gormConfig.goVersion = runtime.Version()

	c.log.gorm = true

	return c, nil
}

// 创建模型
func (c *ApiClient) gormAutoMigrate() error {
	return c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&apiPostgresqlLog{})
}

// 模型结构体
type apiPostgresqlLog struct {
	LogId                 uint           `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`                //【记录】编号
	TraceId               string         `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`                 //【系统】跟踪编号
	RequestTime           time.Time      `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`               //【请求】时间
	RequestUri            string         `gorm:"comment:【请求】链接" json:"request_uri,omitempty"`                      //【请求】链接
	RequestUrl            string         `gorm:"comment:【请求】链接" json:"request_url,omitempty"`                      //【请求】链接
	RequestApi            string         `gorm:"index;comment:【请求】接口" json:"request_api,omitempty"`                //【请求】接口
	RequestMethod         string         `gorm:"index;comment:【请求】方式" json:"request_method,omitempty"`             //【请求】方式
	RequestParams         datatypes.JSON `gorm:"type:jsonb;comment:【请求】参数" json:"request_params,omitempty"`        //【请求】参数
	RequestHeader         datatypes.JSON `gorm:"type:jsonb;comment:【请求】头部" json:"request_header,omitempty"`        //【请求】头部
	ResponseHeader        datatypes.JSON `gorm:"type:jsonb;comment:【返回】头部" json:"response_header,omitempty"`       //【返回】头部
	ResponseStatusCode    int            `gorm:"index;comment:【返回】状态码" json:"response_status_code,omitempty"`      //【返回】状态码
	ResponseBody          datatypes.JSON `gorm:"type:jsonb;comment:【返回】内容" json:"response_body,omitempty"`         //【返回】内容
	ResponseBodyXml       string         `gorm:"type:xml;comment:【返回】内容 Xml格式" json:"response_body_xml,omitempty"` //【返回】内容 Xml格式
	ResponseContentLength int64          `gorm:"comment:【返回】大小" json:"response_content_length,omitempty"`          //【返回】大小
	ResponseTime          time.Time      `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`              //【返回】时间
	SystemHostName        string         `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"`          //【系统】主机名
	SystemInsideIp        string         `gorm:"index;comment:【系统】内网ip" json:"system_inside_ip,omitempty"`         //【系统】内网ip
	GoVersion             string         `gorm:"index;comment:【程序】Go版本" json:"go_version,omitempty"`               //【程序】Go版本
	SdkVersion            string         `gorm:"index;comment:【程序】Sdk版本" json:"sdk_version,omitempty"`             //【程序】Sdk版本
}

// 记录日志
func (c *ApiClient) gormRecord(ctx context.Context, postgresqlLog apiPostgresqlLog) error {

	if utf8.ValidString(string(postgresqlLog.ResponseBody)) == false {
		postgresqlLog.ResponseBody = datatypes.JSON("")
	}

	postgresqlLog.SystemHostName = c.gormConfig.hostname
	if postgresqlLog.SystemInsideIp == "" {
		postgresqlLog.SystemInsideIp = c.gormConfig.insideIp
	}
	postgresqlLog.GoVersion = c.gormConfig.goVersion

	postgresqlLog.TraceId = gotrace_id.GetTraceIdContext(ctx)

	return c.gormClient.Db.Table(c.gormConfig.tableName).Create(&postgresqlLog).Error
}

// GormQuery 查询
func (c *ApiClient) GormQuery() *gorm.DB {
	return c.gormClient.Db.Table(c.gormConfig.tableName)
}

// GormMiddleware 中间件
func (c *ApiClient) GormMiddleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiPostgresqlLog{
		RequestTime:           request.RequestTime,                                            //【请求】时间
		RequestUri:            request.RequestUri,                                             //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,                         //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,                        //【请求】接口
		RequestMethod:         request.RequestMethod,                                          //【请求】方式
		RequestParams:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestParams)),  //【请求】参数
		RequestHeader:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestHeader)),  //【请求】头部
		ResponseHeader:        datatypes.JSON(dorm.JsonEncodeNoError(request.ResponseHeader)), //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                                     //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                                  //【返回】大小
		ResponseTime:          request.ResponseTime,                                           //【返回】时间
		SdkVersion:            sdkVersion,                                                     //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[log.GormMiddleware]：%s %s\n", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(dorm.JsonDecodeNoError(request.ResponseBody)) > 0 {
			data.ResponseBody = request.ResponseBody //【返回】内容
		} else {
			c.zapLog.WithTraceId(ctx).Sugar().Infof("[log.GormMiddleware]：%s %s\n", data.RequestUri, request.ResponseBody)
		}
	}
	err := c.gormRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[log.GormMiddleware]：%s\n", err.Error())
	}
}

// GormMiddlewareXml 中间件
func (c *ApiClient) GormMiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiPostgresqlLog{
		RequestTime:           request.RequestTime,                                            //【请求】时间
		RequestUri:            request.RequestUri,                                             //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,                         //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,                        //【请求】接口
		RequestMethod:         request.RequestMethod,                                          //【请求】方式
		RequestParams:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestParams)),  //【请求】参数
		RequestHeader:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestHeader)),  //【请求】头部
		ResponseHeader:        datatypes.JSON(dorm.JsonEncodeNoError(request.ResponseHeader)), //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                                     //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                                  //【返回】大小
		ResponseTime:          request.ResponseTime,                                           //【返回】时间
		SdkVersion:            sdkVersion,                                                     //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[log.GormMiddlewareXml]：%s %s\n", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(string(request.ResponseBody)) > 0 {
			data.ResponseBodyXml = string(request.ResponseBody) //【返回】内容  Xml格式
		} else {
			c.zapLog.WithTraceId(ctx).Sugar().Infof("[log.GormMiddlewareXml]：%s %s\n", data.RequestUri, request.ResponseBody)
		}
	}
	err := c.gormRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[log.GormMiddlewareXml]：%s\n", err.Error())
	}
}

// GormMiddlewareCustom 中间件
func (c *ApiClient) GormMiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	data := apiPostgresqlLog{
		RequestTime:           request.RequestTime,                                            //【请求】时间
		RequestUri:            request.RequestUri,                                             //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,                         //【请求】链接
		RequestApi:            api,                                                            //【请求】接口
		RequestMethod:         request.RequestMethod,                                          //【请求】方式
		RequestParams:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestParams)),  //【请求】参数
		RequestHeader:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestHeader)),  //【请求】头部
		ResponseHeader:        datatypes.JSON(dorm.JsonEncodeNoError(request.ResponseHeader)), //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                                     //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                                  //【返回】大小
		ResponseTime:          request.ResponseTime,                                           //【返回】时间
		SdkVersion:            sdkVersion,                                                     //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[log.GormMiddlewareCustom]：%s %s\n", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(dorm.JsonDecodeNoError(request.ResponseBody)) > 0 {
			data.ResponseBody = request.ResponseBody //【返回】内容
		} else {
			c.zapLog.WithTraceId(ctx).Sugar().Infof("[log.GormMiddlewareCustom]：%s %s\n", data.RequestUri, request.ResponseBody)
		}
	}
	err := c.gormRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[log.GormMiddlewareCustom]：%s\n", err.Error())
	}
}
