package golog

import (
	"context"
	"go.dtapp.net/dorm"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 模型结构体
type apiMongolLog struct {
	LogId                 primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                                      //【记录】编号
	LogTime               primitive.DateTime `json:"log_time,omitempty" bson:"log_time"`                                         //【记录】时间
	TraceId               string             `json:"trace_id,omitempty" bson:"trace_id,omitempty"`                               //【记录】跟踪编号
	RequestTime           dorm.BsonTime      `json:"request_time,omitempty" bson:"request_time,omitempty"`                       //【请求】时间
	RequestUri            string             `json:"request_uri,omitempty" bson:"request_uri,omitempty"`                         //【请求】链接
	RequestUrl            string             `json:"request_url,omitempty" bson:"request_url,omitempty"`                         //【请求】链接
	RequestApi            string             `json:"request_api,omitempty" bson:"request_api,omitempty"`                         //【请求】接口
	RequestMethod         string             `json:"request_method,omitempty" bson:"request_method,omitempty"`                   //【请求】方式
	RequestParams         interface{}        `json:"request_params,omitempty" bson:"request_params,omitempty"`                   //【请求】参数
	RequestHeader         interface{}        `json:"request_header,omitempty" bson:"request_header,omitempty"`                   //【请求】头部
	RequestIp             string             `json:"request_ip,omitempty" bson:"request_ip,omitempty"`                           //【请求】请求Ip
	ResponseHeader        interface{}        `json:"response_header,omitempty" bson:"response_header,omitempty"`                 //【返回】头部
	ResponseStatusCode    int                `json:"response_status_code,omitempty" bson:"response_status_code,omitempty"`       //【返回】状态码
	ResponseBody          interface{}        `json:"response_body,omitempty" bson:"response_body,omitempty"`                     //【返回】内容
	ResponseContentLength int64              `json:"response_content_length,omitempty" bson:"response_content_length,omitempty"` //【返回】大小
	ResponseTime          dorm.BsonTime      `json:"response_time,omitempty" bson:"response_time,omitempty"`                     //【返回】时间
	SystemHostName        string             `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`               //【系统】主机名
	SystemInsideIp        string             `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`               //【系统】内网ip
	SystemOs              string             `json:"system_os,omitempty" bson:"system_os,omitempty"`                             //【系统】系统类型
	SystemArch            string             `json:"system_arch,omitempty" bson:"system_arch,omitempty"`                         //【系统】系统架构
	GoVersion             string             `json:"go_version,omitempty" bson:"go_version,omitempty"`                           //【程序】Go版本
	SdkVersion            string             `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`                         //【程序】Sdk版本
}

// 创建时间序列集合
func (c *ApiClient) mongoCreateCollection(ctx context.Context) {
	err := c.mongoClient.Database(c.mongoConfig.databaseName).CreateCollection(ctx, c.mongoConfig.collectionName, options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField("log_time")))
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建时间序列集合：%s", err)
	}
}

// 创建索引
func (c *ApiClient) mongoCreateIndexes(ctx context.Context) {
	indexes, err := c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).CreateManyIndexes(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{
				Key:   "log_time",
				Value: -1,
			}},
		}})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建索引：%s", err)
	}
	c.zapLog.WithTraceId(ctx).Sugar().Infof("创建索引：%s", indexes)
}

// 记录日志
func (c *ApiClient) mongoRecord(ctx context.Context, mongoLog apiMongolLog) (err error) {

	mongoLog.SystemHostName = c.config.systemHostName    //【系统】主机名
	mongoLog.SystemInsideIp = c.config.systemInsideIp    //【系统】内网ip
	mongoLog.GoVersion = c.config.goVersion              //【程序】Go版本
	mongoLog.TraceId = gotrace_id.GetTraceIdContext(ctx) //【记录】跟踪编号
	mongoLog.RequestIp = c.config.systemOutsideIp        //【请求】请求Ip
	mongoLog.SystemOs = c.config.systemOs                //【系统】系统类型
	mongoLog.SystemArch = c.config.systemArch            //【系统】系统架构
	mongoLog.LogId = primitive.NewObjectID()             //【记录】编号

	_, err = c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).InsertOne(ctx, mongoLog)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("记录日志失败：%s", err)
	}
	return err
}

// MongoDelete 删除
func (c *ApiClient) MongoDelete(ctx context.Context, hour int64) (*mongo.DeleteResult, error) {
	filter := bson.D{{"log_time", bson.D{{"$lt", primitive.NewDateTimeFromTime(gotime.Current().BeforeHour(hour).Time)}}}}
	return c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).DeleteMany(ctx, filter)
}

// 中间件
func (c *ApiClient) mongoMiddleware(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiMongolLog{
		LogTime:               primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:           dorm.NewBsonTimeFromTime(request.RequestTime),      //【请求】时间
		RequestUri:            request.RequestUri,                                 //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,            //【请求】接口
		RequestMethod:         request.RequestMethod,                              //【请求】方式
		RequestParams:         request.RequestParams,                              //【请求】参数
		RequestHeader:         request.RequestHeader,                              //【请求】头部
		ResponseHeader:        request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                         //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                      //【返回】大小
		ResponseTime:          dorm.NewBsonTimeFromTime(request.ResponseTime),     //【返回】时间
		SdkVersion:            sdkVersion,                                         //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddleware.type]：%s %s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.JsonDecodeNoError(request.ResponseBody) //【返回】内容
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddleware.len]：%s %s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddleware.data]：%+v", data)
	}

	err := c.mongoRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.mongoMiddleware]：%s", err.Error())
	}
}

// 中间件
func (c *ApiClient) mongoMiddlewareXml(ctx context.Context, request gorequest.Response, sdkVersion string) {
	data := apiMongolLog{
		LogTime:               primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:           dorm.NewBsonTimeFromTime(request.RequestTime),      //【请求】时间
		RequestUri:            request.RequestUri,                                 //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:            gourl.UriParse(request.RequestUri).Path,            //【请求】接口
		RequestMethod:         request.RequestMethod,                              //【请求】方式
		RequestParams:         request.RequestParams,                              //【请求】参数
		RequestHeader:         request.RequestHeader,                              //【请求】头部
		ResponseHeader:        request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                         //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                      //【返回】大小
		ResponseTime:          dorm.NewBsonTimeFromTime(request.ResponseTime),     //【返回】时间
		SdkVersion:            sdkVersion,                                         //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddlewareXml.type]：%s %s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.XmlDecodeNoError(request.ResponseBody) //【返回】内容
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddlewareXml]：%s %s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddlewareXml.data]：%+v", data)
	}

	err := c.mongoRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.mongoMiddlewareXml]：%s", err.Error())
	}
}

// 中间件
func (c *ApiClient) mongoMiddlewareCustom(ctx context.Context, api string, request gorequest.Response, sdkVersion string) {
	data := apiMongolLog{
		LogTime:               primitive.NewDateTimeFromTime(request.RequestTime), //【记录】时间
		RequestTime:           dorm.NewBsonTimeFromTime(request.RequestTime),      //【请求】时间
		RequestUri:            request.RequestUri,                                 //【请求】链接
		RequestUrl:            gourl.UriParse(request.RequestUri).Url,             //【请求】链接
		RequestApi:            api,                                                //【请求】接口
		RequestMethod:         request.RequestMethod,                              //【请求】方式
		RequestParams:         request.RequestParams,                              //【请求】参数
		RequestHeader:         request.RequestHeader,                              //【请求】头部
		ResponseHeader:        request.ResponseHeader,                             //【返回】头部
		ResponseStatusCode:    request.ResponseStatusCode,                         //【返回】状态码
		ResponseContentLength: request.ResponseContentLength,                      //【返回】大小
		ResponseTime:          dorm.NewBsonTimeFromTime(request.ResponseTime),     //【返回】时间
		SdkVersion:            sdkVersion,                                         //【程序】Sdk版本
	}
	if request.HeaderIsImg() {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddlewareCustom.type]：%s %s", data.RequestUri, request.ResponseHeader.Get("Content-Type"))
	} else {
		if len(request.ResponseBody) > 0 {
			data.ResponseBody = dorm.JsonDecodeNoError(request.ResponseBody) //【返回】内容
		} else {
			if c.logDebug {
				c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddlewareCustom]：%s %s", data.RequestUri, request.ResponseBody)
			}
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceId(ctx).Sugar().Infof("[golog.api.mongoMiddlewareCustom.data]：%+v", data)
	}

	err := c.mongoRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[golog.api.mongoMiddlewareCustom]：%s", err.Error())
	}
}
