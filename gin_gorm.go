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
	"go.dtapp.net/gourl"
	"go.dtapp.net/goxml"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"time"
)

// GinGormClientConfig 框架实例配置
type GinGormClientConfig struct {
	IpService     *goip.Client     // ip服务
	GormClientFun ginGormClientFun // 日志配置
	Debug         bool             // 日志开关
	ZapLog        *ZapLog          // 日志服务
}

// NewGinGormClient 创建框架实例化
// client 数据库服务
// tableName 表名
// ipService ip服务
func NewGinGormClient(config *GinGormClientConfig) (*GinClient, error) {

	var ctx = context.Background()

	c := &GinClient{}

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

	c.ipService = config.IpService

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
func (c *GinClient) gormAutoMigrate() error {
	return c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&ginPostgresqlLog{})
}

// 模型结构体
type ginPostgresqlLog struct {
	LogId             uint           `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`                 //【记录】编号
	TraceId           string         `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`                  //【系统】跟踪编号
	RequestTime       time.Time      `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`                //【请求】时间
	RequestUri        string         `gorm:"comment:【请求】请求链接 域名+路径+参数" json:"request_uri,omitempty"`            //【请求】请求链接 域名+路径+参数
	RequestUrl        string         `gorm:"comment:【请求】请求链接 域名+路径" json:"request_url,omitempty"`               //【请求】请求链接 域名+路径
	RequestApi        string         `gorm:"index;comment:【请求】请求接口 路径" json:"request_api,omitempty"`            //【请求】请求接口 路径
	RequestMethod     string         `gorm:"index;comment:【请求】请求方式" json:"request_method,omitempty"`            //【请求】请求方式
	RequestProto      string         `gorm:"comment:【请求】请求协议" json:"request_proto,omitempty"`                   //【请求】请求协议
	RequestUa         string         `gorm:"comment:【请求】请求UA" json:"request_ua,omitempty"`                      //【请求】请求UA
	RequestReferer    string         `gorm:"comment:【请求】请求referer" json:"request_referer,omitempty"`            //【请求】请求referer
	RequestBody       datatypes.JSON `gorm:"type:jsonb;comment:【请求】请求主体" json:"request_body,omitempty"`         //【请求】请求主体
	RequestBodyXml    string         `gorm:"type:xml;comment:【请求】请求主体 Xml格式" json:"request_body_xml,omitempty"` //【请求】请求主体 Xml格式
	RequestUrlQuery   datatypes.JSON `gorm:"type:jsonb;comment:【请求】请求URL参数" json:"request_url_query,omitempty"` //【请求】请求URL参数
	RequestIp         string         `gorm:"index;comment:【请求】请求客户端Ip" json:"request_ip,omitempty"`             //【请求】请求客户端Ip
	RequestIpCountry  string         `gorm:"index;comment:【请求】请求客户端城市" json:"request_ip_country,omitempty"`     //【请求】请求客户端城市
	RequestIpRegion   string         `gorm:"index;comment:【请求】请求客户端区域" json:"request_ip_region,omitempty"`      //【请求】请求客户端区域
	RequestIpProvince string         `gorm:"index;comment:【请求】请求客户端省份" json:"request_ip_province,omitempty"`    //【请求】请求客户端省份
	RequestIpCity     string         `gorm:"index;comment:【请求】请求客户端城市" json:"request_ip_city,omitempty"`        //【请求】请求客户端城市
	RequestIpIsp      string         `gorm:"index;comment:【请求】请求客户端运营商" json:"request_ip_isp,omitempty"`        //【请求】请求客户端运营商
	RequestHeader     datatypes.JSON `gorm:"type:jsonb;comment:【请求】请求头" json:"request_header,omitempty"`        //【请求】请求头
	ResponseTime      time.Time      `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`               //【返回】时间
	ResponseCode      int            `gorm:"index;comment:【返回】状态码" json:"response_code,omitempty"`              //【返回】状态码
	ResponseMsg       string         `gorm:"comment:【返回】描述" json:"response_msg,omitempty"`                      //【返回】描述
	ResponseData      datatypes.JSON `gorm:"type:jsonb;comment:【返回】数据" json:"response_data,omitempty"`          //【返回】数据
	CostTime          int64          `gorm:"comment:【系统】花费时间" json:"cost_time,omitempty"`                       //【系统】花费时间
	SystemHostName    string         `gorm:"index;comment:【系统】主机名" json:"system_host_name,omitempty"`           //【系统】主机名
	SystemInsideIp    string         `gorm:"index;comment:【系统】内网ip" json:"system_inside_ip,omitempty"`          //【系统】内网ip
	GoVersion         string         `gorm:"index;comment:【程序】Go版本" json:"go_version,omitempty"`                //【程序】Go版本
	SdkVersion        string         `gorm:"index;comment:【程序】Sdk版本" json:"sdk_version,omitempty"`              //【程序】Sdk版本
}

// gormRecord 记录日志
func (c *GinClient) gormRecord(postgresqlLog ginPostgresqlLog) error {

	postgresqlLog.SystemHostName = c.gormConfig.hostname
	if postgresqlLog.SystemInsideIp == "" {
		postgresqlLog.SystemInsideIp = c.gormConfig.insideIp
	}
	postgresqlLog.GoVersion = c.gormConfig.goVersion

	postgresqlLog.SdkVersion = Version

	return c.gormClient.Db.Table(c.gormConfig.tableName).Create(&postgresqlLog).Error
}

func (c *GinClient) gormRecordJson(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody map[string]interface{}, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string) {
	data := ginPostgresqlLog{
		TraceId:           traceId,                                                            //【系统】跟踪编号
		RequestTime:       requestTime,                                                        //【请求】时间
		RequestUrl:        ginCtx.Request.RequestURI,                                          //【请求】请求链接
		RequestApi:        gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI),       //【请求】请求接口
		RequestMethod:     ginCtx.Request.Method,                                              //【请求】请求方式
		RequestProto:      ginCtx.Request.Proto,                                               //【请求】请求协议
		RequestUa:         ginCtx.Request.UserAgent(),                                         //【请求】请求UA
		RequestReferer:    ginCtx.Request.Referer(),                                           //【请求】请求referer
		RequestUrlQuery:   datatypes.JSON(dorm.JsonEncodeNoError(ginCtx.Request.URL.Query())), //【请求】请求URL参数
		RequestIp:         clientIp,                                                           //【请求】请求客户端Ip
		RequestIpCountry:  requestClientIpCountry,                                             //【请求】请求客户端城市
		RequestIpRegion:   requestClientIpRegion,                                              //【请求】请求客户端区域
		RequestIpProvince: requestClientIpProvince,                                            //【请求】请求客户端省份
		RequestIpCity:     requestClientIpCity,                                                //【请求】请求客户端城市
		RequestIpIsp:      requestClientIpIsp,                                                 //【请求】请求客户端运营商
		RequestHeader:     datatypes.JSON(dorm.JsonEncodeNoError(ginCtx.Request.Header)),      //【请求】请求头
		ResponseTime:      gotime.Current().Time,                                              //【返回】时间
		ResponseCode:      responseCode,                                                       //【返回】状态码
		ResponseData:      datatypes.JSON(responseBody),                                       //【返回】数据
		CostTime:          endTime - startTime,                                                //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(dorm.JsonEncodeNoError(requestBody)) > 0 {
		data.RequestBody = datatypes.JSON(dorm.JsonEncodeNoError(requestBody)) //【请求】请求主体
	} else {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[log.gormRecordJson]：%s %s\n", data.RequestUri, requestBody)
	}

	err := c.gormRecord(data)
	if err != nil {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gormRecordJson]：%s\n", err)
	}
}

func (c *GinClient) gormRecordXml(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody string, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string) {
	data := ginPostgresqlLog{
		TraceId:           traceId,                                                            //【系统】跟踪编号
		RequestTime:       requestTime,                                                        //【请求】时间
		RequestUrl:        ginCtx.Request.RequestURI,                                          //【请求】请求链接
		RequestApi:        gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI),       //【请求】请求接口
		RequestMethod:     ginCtx.Request.Method,                                              //【请求】请求方式
		RequestProto:      ginCtx.Request.Proto,                                               //【请求】请求协议
		RequestUa:         ginCtx.Request.UserAgent(),                                         //【请求】请求UA
		RequestReferer:    ginCtx.Request.Referer(),                                           //【请求】请求referer
		RequestUrlQuery:   datatypes.JSON(dorm.JsonEncodeNoError(ginCtx.Request.URL.Query())), //【请求】请求URL参数
		RequestIp:         clientIp,                                                           //【请求】请求客户端Ip
		RequestIpCountry:  requestClientIpCountry,                                             //【请求】请求客户端城市
		RequestIpRegion:   requestClientIpRegion,                                              //【请求】请求客户端区域
		RequestIpProvince: requestClientIpProvince,                                            //【请求】请求客户端省份
		RequestIpCity:     requestClientIpCity,                                                //【请求】请求客户端城市
		RequestIpIsp:      requestClientIpIsp,                                                 //【请求】请求客户端运营商
		RequestHeader:     datatypes.JSON(dorm.JsonEncodeNoError(ginCtx.Request.Header)),      //【请求】请求头
		ResponseTime:      gotime.Current().Time,                                              //【返回】时间
		ResponseCode:      responseCode,                                                       //【返回】状态码
		ResponseData:      datatypes.JSON(responseBody),                                       //【返回】数据
		CostTime:          endTime - startTime,                                                //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBodyXml = requestBody //【请求】请求主体  Xml格式
	} else {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[log.gormRecordXml]：%s %s\n", data.RequestUri, requestBody)
	}

	err := c.gormRecord(data)
	if err != nil {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gormRecordXml]：%s\n", err)
	}
}

// GormQuery 查询
func (c *GinClient) GormQuery() *gorm.DB {
	return c.gormClient.Db.Table(c.gormConfig.tableName)
}

// GormMiddleware 中间件
func (c *GinClient) GormMiddleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		// 开始时间
		startTime := gotime.Current().TimestampWithMillisecond()
		requestTime := gotime.Current().Time

		// 获取
		data, _ := ioutil.ReadAll(ginCtx.Request.Body)

		if c.gormConfig.debug {
			c.zapLog.WithLogger().Sugar().Infof("[golog.GormMiddleware] %s\n", data)
		}

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

				if c.gormConfig.debug {
					c.zapLog.WithLogger().Sugar().Infof("[golog.GormMiddleware.len(jsonBody)] %v\n", len(jsonBody))
				}

				if err != nil {
					if c.gormConfig.debug {
						c.zapLog.WithLogger().Sugar().Infof("[golog.GormMiddleware.json.Unmarshal] %s %s\n", jsonBody, err)
					}
					dataJson = false
					xmlBody = goxml.XmlDecode(string(data))
				}
			}

			if c.gormConfig.debug {
				c.zapLog.WithLogger().Sugar().Infof("[golog.GormMiddleware.xmlBody] %s\n", xmlBody)
				c.zapLog.WithLogger().Sugar().Infof("[golog.GormMiddleware.jsonBody] %s\n", jsonBody)
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

			// 记录
			if c.gormClient != nil && c.gormClient.Db != nil {

				var traceId = gotrace_id.GetGinTraceId(ginCtx)

				if dataJson {
					if c.gormConfig.debug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.GormMiddleware.gormRecord.json.request_body] %s\n", jsonBody)
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.GormMiddleware.gormRecord.json.request_body] %s\n", dorm.JsonEncodeNoError(jsonBody))
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.GormMiddleware.gormRecord.json.request_body] %s\n", datatypes.JSON(dorm.JsonEncodeNoError(jsonBody)))
					}
					c.gormRecordJson(ginCtx, traceId, requestTime, jsonBody, responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				} else {
					if c.gormConfig.debug {
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.GormMiddleware.gormRecord.xml.request_body] %s\n", xmlBody)
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.GormMiddleware.gormRecord.xml.request_body] %s\n", dorm.JsonEncodeNoError(xmlBody))
						c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.GormMiddleware.gormRecord.xml.request_body] %s\n", datatypes.JSON(dorm.JsonEncodeNoError(xmlBody)))
					}
					c.gormRecordXml(ginCtx, traceId, requestTime, string(data), responseCode, responseBody, startTime, endTime, clientIp, requestClientIpCountry, requestClientIpRegion, requestClientIpProvince, requestClientIpCity, requestClientIpIsp)
				}
			}
		}()
	}
}
