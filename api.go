package golog

import (
	"context"
	"errors"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gourl"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/datatypes"
	"log"
	"os"
	"runtime"
)

// ApiClient 接口
type ApiClient struct {
	gormClient  *dorm.GormClient  // 数据库驱动
	mongoClient *dorm.MongoClient // 数据库驱动
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
type apiGormClientFun func() (*dorm.GormClient, string)

// client 数据库服务
// string 库名
// string 表名
type apiMongoClientFun func() (*dorm.MongoClient, string, string)

// ApiClientConfig 接口实例配置
type ApiClientConfig struct {
	GormClientFun  apiGormClientFun  // 日志配置
	MongoClientFun apiMongoClientFun // 日志配置
	Debug          bool              // 日志开关
}

// NewApiClient 创建接口实例化
// client 数据库服务
// tableName 表名
func NewApiClient(config *ApiClientConfig) (*ApiClient, error) {

	var ctx = context.Background()

	c := &ApiClient{}

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

		err := c.gormClient.Db.Table(c.gormConfig.tableName).AutoMigrate(&apiPostgresqlLog{})
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

		c.mongoConfig.hostname = hostname
		c.mongoConfig.insideIp = goip.GetInsideIp(ctx)
		c.mongoConfig.goVersion = runtime.Version()

		c.log.mongo = true
	}

	return c, nil
}

// Middleware 中间件
func (c *ApiClient) Middleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		request.ResponseBody = []byte{}
	}
	if c.log.gorm {
		err := c.gormRecord(ctx, apiPostgresqlLog{
			RequestTime:           request.RequestTime,                                            //【请求】时间
			RequestUri:            request.RequestUri,                                             //【请求】链接
			RequestUrl:            gourl.UriParse(request.RequestUri).Url,                         //【请求】链接
			RequestApi:            gourl.UriParse(request.RequestUri).Path,                        //【请求】接口
			RequestMethod:         request.RequestMethod,                                          //【请求】方式
			RequestParams:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestParams)),  //【请求】参数
			RequestHeader:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestHeader)),  //【请求】头部
			ResponseHeader:        datatypes.JSON(dorm.JsonEncodeNoError(request.ResponseHeader)), //【返回】头部
			ResponseStatusCode:    request.ResponseStatusCode,                                     //【返回】状态码
			ResponseBody:          request.ResponseBody,                                           //【返回】内容
			ResponseContentLength: request.ResponseContentLength,                                  //【返回】大小
			ResponseTime:          request.ResponseTime,                                           //【返回】时间
			SdkVersion:            sdkVersion,                                                     //【程序】Sdk版本
		})
		if err != nil {
			if c.gormConfig.debug {
				log.Printf("[log.Middleware]%s\n", err.Error())
			}
		}
	}
	if c.log.mongo {
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
				log.Printf("[log.Middleware]%s\n", err.Error())
			}
		}
	}
}

// MiddlewareXml 中间件
func (c *ApiClient) MiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		request.ResponseBody = []byte{}
	}
	if c.log.gorm {
		err := c.gormRecord(ctx, apiPostgresqlLog{
			RequestTime:           request.RequestTime,                                                                 //【请求】时间
			RequestUri:            request.RequestUri,                                                                  //【请求】链接
			RequestUrl:            gourl.UriParse(request.RequestUri).Url,                                              //【请求】链接
			RequestApi:            gourl.UriParse(request.RequestUri).Path,                                             //【请求】接口
			RequestMethod:         request.RequestMethod,                                                               //【请求】方式
			RequestParams:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestParams)),                       //【请求】参数
			RequestHeader:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestHeader)),                       //【请求】头部
			ResponseHeader:        datatypes.JSON(dorm.JsonEncodeNoError(request.ResponseHeader)),                      //【返回】头部
			ResponseStatusCode:    request.ResponseStatusCode,                                                          //【返回】状态码
			ResponseBody:          datatypes.JSON(dorm.JsonEncodeNoError(dorm.XmlDecodeNoError(request.ResponseBody))), //【返回】内容
			ResponseContentLength: request.ResponseContentLength,                                                       //【返回】大小
			ResponseTime:          request.ResponseTime,                                                                //【返回】时间
			SdkVersion:            sdkVersion,                                                                          //【程序】Sdk版本
		})
		if err != nil {
			if c.gormConfig.debug {
				log.Printf("[log.MiddlewareXml]%s\n", err.Error())
			}
		}
	}
	if c.log.mongo {
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
				log.Printf("[log.MiddlewareXml]%s\n", err.Error())
			}
		}
	}
}

// MiddlewareCustom 中间件
func (c *ApiClient) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	if request.ResponseHeader.Get("Content-Type") == "image/jpeg" || request.ResponseHeader.Get("Content-Type") == "image/png" || request.ResponseHeader.Get("Content-Type") == "image/jpg" {
		request.ResponseBody = []byte{}
	}
	if c.log.gorm {
		err := c.gormRecord(ctx, apiPostgresqlLog{
			RequestTime:           request.RequestTime,                                            //【请求】时间
			RequestUri:            request.RequestUri,                                             //【请求】链接
			RequestUrl:            gourl.UriParse(request.RequestUri).Url,                         //【请求】链接
			RequestApi:            api,                                                            //【请求】接口
			RequestMethod:         request.RequestMethod,                                          //【请求】方式
			RequestParams:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestParams)),  //【请求】参数
			RequestHeader:         datatypes.JSON(dorm.JsonEncodeNoError(request.RequestHeader)),  //【请求】头部
			ResponseHeader:        datatypes.JSON(dorm.JsonEncodeNoError(request.ResponseHeader)), //【返回】头部
			ResponseStatusCode:    request.ResponseStatusCode,                                     //【返回】状态码
			ResponseBody:          request.ResponseBody,                                           //【返回】内容
			ResponseContentLength: request.ResponseContentLength,                                  //【返回】大小
			ResponseTime:          request.ResponseTime,                                           //【返回】时间
			SdkVersion:            sdkVersion,                                                     //【程序】Sdk版本
		})
		if err != nil {
			if c.gormConfig.debug {
				log.Printf("[log.MiddlewareCustom]%s\n", err.Error())
			}
		}
	}
	if c.log.mongo {
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
				log.Printf("[log.MiddlewareCustom]%s\n", err.Error())
			}
		}
	}
}
