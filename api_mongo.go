package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"runtime"
)

// ApiMongoClientConfig 接口实例配置
type ApiMongoClientConfig struct {
	MongoClientFun apiMongoClientFun // 日志配置
	Debug          bool              // 日志开关
}

// NewApiMongoClient 创建接口实例化
// client 数据库服务
// databaseName 库名
// collectionName 表名
func NewApiMongoClient(config *ApiMongoClientConfig) (*ApiClient, error) {

	var ctx = context.Background()

	c := &ApiClient{}

	client, databaseName, collectionName := config.MongoClientFun()

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

	c.mongoConfig.debug = config.Debug

	hostname, _ := os.Hostname()

	c.mongoConfig.hostname = hostname
	c.mongoConfig.insideIp = goip.GetInsideIp(ctx)
	c.mongoConfig.goVersion = runtime.Version()

	return c, nil
}

// 模型结构体
type apiMongolLog struct {
	LogId                 primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                                      //【记录】编号
	TraceId               string             `json:"trace_id,omitempty" bson:"trace_id,omitempty"`                               //【系统】跟踪编号
	RequestTime           primitive.DateTime `json:"request_time,omitempty" bson:"request_time,omitempty"`                       //【请求】时间
	RequestUri            string             `json:"request_uri,omitempty" bson:"request_uri,omitempty"`                         //【请求】链接
	RequestUrl            string             `json:"request_url,omitempty" bson:"request_url,omitempty"`                         //【请求】链接
	RequestApi            string             `json:"request_api,omitempty" bson:"request_api,omitempty"`                         //【请求】接口
	RequestMethod         string             `json:"request_method,omitempty" bson:"request_method,omitempty"`                   //【请求】方式
	RequestParams         interface{}        `json:"request_params,omitempty" bson:"request_params,omitempty"`                   //【请求】参数
	RequestHeader         interface{}        `json:"request_header,omitempty" bson:"request_header,omitempty"`                   //【请求】头部
	ResponseHeader        interface{}        `json:"response_header,omitempty" bson:"response_header,omitempty"`                 //【返回】头部
	ResponseStatusCode    int                `json:"response_status_code,omitempty" bson:"response_status_code,omitempty"`       //【返回】状态码
	ResponseBody          interface{}        `json:"response_body,omitempty" bson:"response_body,omitempty"`                     //【返回】内容
	ResponseContentLength int64              `json:"response_content_length,omitempty" bson:"response_content_length,omitempty"` //【返回】大小
	ResponseTime          primitive.DateTime `json:"response_time,omitempty" bson:"response_time,omitempty"`                     //【返回】时间
	SystemHostName        string             `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`               //【系统】主机名
	SystemInsideIp        string             `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`               //【系统】内网ip
	GoVersion             string             `json:"go_version,omitempty" bson:"go_version,omitempty"`                           //【程序】Go版本
	SdkVersion            string             `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`                         //【程序】Sdk版本
}

// 记录日志
func (c *ApiClient) mongoRecord(ctx context.Context, mongoLog apiMongolLog) error {

	mongoLog.SystemHostName = c.mongoConfig.hostname
	if mongoLog.SystemInsideIp == "" {
		mongoLog.SystemInsideIp = c.mongoConfig.insideIp
	}
	mongoLog.GoVersion = c.mongoConfig.goVersion

	mongoLog.TraceId = gotrace_id.GetTraceIdContext(ctx)

	mongoLog.LogId = primitive.NewObjectID()

	_, err := c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).InsertOne(mongoLog)

	return err
}

// MongoQuery 查询
func (c *ApiClient) MongoQuery() *dorm.MongoClient {
	return c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName)
}

// MongoMiddleware 中间件
func (c *ApiClient) MongoMiddleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		request.ResponseBody = []byte{}
	}
	err := c.mongoRecord(ctx, apiMongolLog{
		RequestTime:           primitive.NewDateTimeFromTime(request.RequestTime),  //【请求】时间
		RequestUri:            request.RequestUri,                                  //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,              //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,             //【请求】接口
		RequestMethod:         request.RequestMethod,                               //【请求】方式
		RequestParams:         request.RequestParams,                               //【请求】参数
		RequestHeader:         request.RequestHeader,                               //【请求】头部
		ResponseHeader:        request.ResponseHeader,                              //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                          //【返回】状态码
		ResponseBody:          dorm.JsonDecodeNoError(request.ResponseBody),        //【返回】内容
		ResponseContentLength: request.ResponseContentLength,                       //【返回】大小
		ResponseTime:          primitive.NewDateTimeFromTime(request.ResponseTime), //【返回】时间
		SdkVersion:            sdkVersion,                                          //【程序】Sdk版本
	})
	if err != nil {
		if c.mongoConfig.debug {
			log.Printf("[log.MongoMiddleware]%s\n", err.Error())
		}
	}
}

// MongoMiddlewareXml 中间件
func (c *ApiClient) MongoMiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		request.ResponseBody = []byte{}
	}
	err := c.mongoRecord(ctx, apiMongolLog{
		RequestTime:           primitive.NewDateTimeFromTime(request.RequestTime),  //【请求】时间
		RequestUri:            request.RequestUri,                                  //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,              //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,             //【请求】接口
		RequestMethod:         request.RequestMethod,                               //【请求】方式
		RequestParams:         request.RequestParams,                               //【请求】参数
		RequestHeader:         request.RequestHeader,                               //【请求】头部
		ResponseHeader:        request.ResponseHeader,                              //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                          //【返回】状态码
		ResponseBody:          dorm.XmlDecodeNoError(request.ResponseBody),         //【返回】内容
		ResponseContentLength: request.ResponseContentLength,                       //【返回】大小
		ResponseTime:          primitive.NewDateTimeFromTime(request.ResponseTime), //【返回】时间
		SdkVersion:            sdkVersion,                                          //【程序】Sdk版本
	})
	if err != nil {
		if c.mongoConfig.debug {
			log.Printf("[log.MongoMiddlewareXml]%s\n", err.Error())
		}
	}
}

// MongoMiddlewareCustom 中间件
func (c *ApiClient) MongoMiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		request.ResponseBody = []byte{}
	}
	err := c.mongoRecord(ctx, apiMongolLog{
		RequestTime:           primitive.NewDateTimeFromTime(request.RequestTime),  //【请求】时间
		RequestUri:            request.RequestUri,                                  //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,              //【请求】链接
		RequestApi:            api,                                                 //【请求】接口
		RequestMethod:         request.RequestMethod,                               //【请求】方式
		RequestParams:         request.RequestParams,                               //【请求】参数
		RequestHeader:         request.RequestHeader,                               //【请求】头部
		ResponseHeader:        request.ResponseHeader,                              //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                          //【返回】状态码
		ResponseBody:          dorm.JsonDecodeNoError(request.ResponseBody),        //【返回】内容
		ResponseContentLength: request.ResponseContentLength,                       //【返回】大小
		ResponseTime:          primitive.NewDateTimeFromTime(request.ResponseTime), //【返回】时间
		SdkVersion:            sdkVersion,                                          //【程序】Sdk版本
	})
	if err != nil {
		if c.mongoConfig.debug {
			log.Printf("[log.MongoMiddlewareCustom]%s\n", err.Error())
		}
	}
}
