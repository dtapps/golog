package golog

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/dorm"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ginMongoLogRequestIpLocationLocation struct {
	Type        string    `json:"type,omitempty" bson:"type,omitempty"`               // GeoJSON类型
	Coordinates []float64 `json:"coordinates,omitempty" bson:"coordinates,omitempty"` // 经度,纬度
}

// 模型结构体
type ginMongoLog struct {
	LogId             primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                              //【记录】编号
	LogTime           primitive.DateTime `json:"log_time,omitempty" bson:"log_time"`                                 //【记录】时间
	TraceId           string             `json:"trace_id,omitempty" bson:"trace_id,omitempty"`                       //【记录】跟踪编号
	RequestTime       dorm.BsonTime      `json:"request_time,omitempty" bson:"request_time,omitempty"`               //【请求】时间
	RequestUri        string             `json:"request_uri,omitempty" bson:"request_uri,omitempty"`                 //【请求】请求链接 域名+路径+参数
	RequestUrl        string             `json:"request_url,omitempty" bson:"request_url,omitempty"`                 //【请求】请求链接 域名+路径
	RequestApi        string             `json:"request_api,omitempty" bson:"request_api,omitempty"`                 //【请求】请求接口 路径
	RequestMethod     string             `json:"request_method,omitempty" bson:"request_method,omitempty"`           //【请求】请求方式
	RequestProto      string             `json:"request_proto,omitempty" bson:"request_proto,omitempty"`             //【请求】请求协议
	RequestUa         string             `json:"request_ua,omitempty" bson:"request_ua,omitempty"`                   //【请求】请求UA
	RequestReferer    string             `json:"request_referer,omitempty" bson:"request_referer,omitempty"`         //【请求】请求referer
	RequestBody       interface{}        `json:"request_body,omitempty" bson:"request_body,omitempty"`               //【请求】请求主体
	RequestUrlQuery   interface{}        `json:"request_url_query,omitempty" bson:"request_url_query,omitempty"`     //【请求】请求URL参数
	RequestIp         string             `json:"request_ip,omitempty" bson:"request_ip,omitempty"`                   //【请求】请求客户端Ip
	RequestIpCountry  string             `json:"request_ip_country,omitempty" bson:"request_ip_country,omitempty"`   //【请求】请求客户端国家
	RequestIpProvince string             `json:"request_ip_province,omitempty" bson:"request_ip_province,omitempty"` //【请求】请求客户端省份
	RequestIpCity     string             `json:"request_ip_city,omitempty" bson:"request_ip_city,omitempty"`         //【请求】请求客户端城市
	RequestIpIsp      string             `json:"request_ip_isp,omitempty" bson:"request_ip_isp,omitempty"`           //【请求】请求客户端运营商
	RequestIpLocation interface{}        `json:"request_ip_location,omitempty" bson:"request_ip_location,omitempty"` //【请求】请求客户端位置
	RequestHeader     interface{}        `json:"request_header,omitempty" bson:"request_header,omitempty"`           //【请求】请求头
	ResponseTime      dorm.BsonTime      `json:"response_time,omitempty" bson:"response_time,omitempty"`             //【返回】时间
	ResponseCode      int                `json:"response_code,omitempty" bson:"response_code,omitempty"`             //【返回】状态码
	ResponseMsg       string             `json:"response_msg,omitempty" bson:"response_msg,omitempty"`               //【返回】描述
	ResponseData      interface{}        `json:"response_data,omitempty" bson:"response_data,omitempty"`             //【返回】数据
	CostTime          int64              `json:"cost_time,omitempty" bson:"cost_time,omitempty"`                     //【系统】花费时间
	SystemHostName    string             `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`       //【系统】主机名
	SystemInsideIp    string             `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`       //【系统】内网ip
	SystemOs          string             `json:"system_os,omitempty" bson:"system_os,omitempty"`                     //【系统】系统类型
	SystemArch        string             `json:"system_arch,omitempty" bson:"system_arch,omitempty"`                 //【系统】系统架构
	GoVersion         string             `json:"go_version,omitempty" bson:"go_version,omitempty"`                   //【程序】Go版本
	SdkVersion        string             `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`                 //【程序】Sdk版本
}

// 创建集合
func (c *GinClient) mongoCreateCollection(ctx context.Context) {
	err := c.mongoClient.Database(c.mongoConfig.databaseName).CreateCollection(ctx, c.mongoConfig.collectionName, options.CreateCollection().SetCollation(&options.Collation{
		Locale:   "request_time",
		Strength: -1,
	}))
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Error("创建集合：", err)
	}
}

// 创建索引
func (c *GinClient) mongoCreateIndexes(ctx context.Context) {
	indexes, err := c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).CreateManyIndexes(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{
				Key:   "trace_id",
				Value: -1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_method",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_ip",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_ip_country",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_ip_province",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_ip_city",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_ip_isp",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "response_time",
				Value: -1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "response_code",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "system_host_name",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "system_os",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "system_arch",
				Value: -1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "go_version",
				Value: -1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "sdk_version",
				Value: -1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_ip_location",
				Value: "2dsphere",
			}},
		},
	})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建索引：%s", err)
	}
	c.zapLog.WithTraceId(ctx).Sugar().Infof("创建索引：%s", indexes)
}

// 记录日志
func (c *GinClient) mongoRecord(ctx context.Context, mongoLog ginMongoLog) (err error) {

	mongoLog.SystemHostName = c.config.systemHostName //【系统】主机名
	mongoLog.SystemInsideIp = c.config.systemInsideIp //【系统】内网ip
	mongoLog.GoVersion = c.config.goVersion           //【程序】Go版本
	mongoLog.SdkVersion = c.config.sdkVersion         //【程序】Sdk版本
	mongoLog.SystemOs = c.config.systemOs             //【系统】系统类型
	mongoLog.SystemArch = c.config.systemArch         //【系统】系统架构
	mongoLog.LogId = primitive.NewObjectID()          //【记录】编号

	_, err = c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).InsertOne(ctx, mongoLog)
	if err != nil {
		c.zapLog.WithTraceIdStr(mongoLog.TraceId).Sugar().Errorf("记录日志失败：%s", err)
	}
	return err
}

func (c *GinClient) mongoRecordJson(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string, requestClientIpLocationLatitude, requestClientIpLocationLongitude float64) {

	var ctx = gotrace_id.SetGinTraceIdContext(context.Background(), ginCtx)

	if c.logDebug {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.mongoRecordJson]收到保存数据要求：%s,%s", c.mongoConfig.databaseName, c.mongoConfig.collectionName)
	}

	data := ginMongoLog{
		TraceId:           traceId,                                                      //【记录】跟踪编号
		LogTime:           primitive.NewDateTimeFromTime(requestTime),                   //【记录】时间
		RequestTime:       dorm.NewBsonTimeFromTime(requestTime),                        //【请求】时间
		RequestUrl:        ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:        gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:     ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:      ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:         ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:    ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:   ginCtx.Request.URL.Query(),                                   //【请求】请求URL参数
		RequestIp:         clientIp,                                                     //【请求】请求客户端Ip
		RequestIpCountry:  requestClientIpCountry,                                       //【请求】请求客户端国家
		RequestIpProvince: requestClientIpProvince,                                      //【请求】请求客户端省份
		RequestIpCity:     requestClientIpCity,                                          //【请求】请求客户端城市
		RequestIpIsp:      requestClientIpIsp,                                           //【请求】请求客户端运营商
		RequestHeader:     ginCtx.Request.Header,                                        //【请求】请求头
		ResponseTime:      dorm.NewBsonTimeCurrent(),                                    //【返回】时间
		ResponseCode:      responseCode,                                                 //【返回】状态码
		ResponseData:      c.jsonUnmarshal(responseBody),                                //【返回】数据
		CostTime:          endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBody = dorm.JsonDecodeNoError(requestBody) //【请求】请求主体
	} else {
		if c.logDebug {
			c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.mongoRecordJson.len]：%s，%s", data.RequestUri, requestBody)
		}
	}

	if requestClientIpLocationLatitude != 0 && requestClientIpLocationLongitude != 0 {
		data.RequestIpLocation = ginMongoLogRequestIpLocationLocation{
			Type:        "Point",
			Coordinates: []float64{requestClientIpLocationLongitude, requestClientIpLocationLatitude},
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.mongoRecordJson.data]：%+v", data)
	}

	err := c.mongoRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.mongoRecordJson]：%s", err)
	}
}

func (c *GinClient) mongoRecordXml(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, clientIp, requestClientIpCountry, requestClientIpProvince, requestClientIpCity, requestClientIpIsp string, requestClientIpLocationLatitude, requestClientIpLocationLongitude float64) {

	var ctx = gotrace_id.SetGinTraceIdContext(context.Background(), ginCtx)

	if c.logDebug {
		c.zapLog.WithLogger().Sugar().Infof("[golog.gin.mongoRecordXml]收到保存数据要求：%s,%s", c.mongoConfig.databaseName, c.mongoConfig.collectionName)
	}

	data := ginMongoLog{
		TraceId:           traceId,                                                      //【记录】跟踪编号
		LogTime:           primitive.NewDateTimeFromTime(requestTime),                   //【记录】时间
		RequestTime:       dorm.NewBsonTimeFromTime(requestTime),                        //【请求】时间
		RequestUrl:        ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:        gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:     ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:      ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:         ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:    ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery:   ginCtx.Request.URL.Query(),                                   //【请求】请求URL参数
		RequestIp:         clientIp,                                                     //【请求】请求客户端Ip
		RequestIpCountry:  requestClientIpCountry,                                       //【请求】请求客户端国家
		RequestIpProvince: requestClientIpProvince,                                      //【请求】请求客户端省份
		RequestIpCity:     requestClientIpCity,                                          //【请求】请求客户端城市
		RequestIpIsp:      requestClientIpIsp,                                           //【请求】请求客户端运营商
		RequestHeader:     ginCtx.Request.Header,                                        //【请求】请求头
		ResponseTime:      dorm.NewBsonTimeCurrent(),                                    //【返回】时间
		ResponseCode:      responseCode,                                                 //【返回】状态码
		ResponseData:      c.jsonUnmarshal(responseBody),                                //【返回】数据
		CostTime:          endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBody = dorm.XmlDecodeNoError(requestBody) //【请求】请求主体
	} else {
		if c.logDebug {
			c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.mongoRecordXml.len]：%s，%s", data.RequestUri, requestBody)
		}
	}

	if requestClientIpLocationLatitude != 0 && requestClientIpLocationLongitude != 0 {
		data.RequestIpLocation = ginMongoLogRequestIpLocationLocation{
			Type:        "Point",
			Coordinates: []float64{requestClientIpLocationLongitude, requestClientIpLocationLatitude},
		}
	}

	if c.logDebug {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Infof("[golog.gin.mongoRecordXml.data]：%+v", data)
	}

	err := c.mongoRecord(ctx, data)
	if err != nil {
		c.zapLog.WithTraceIdStr(traceId).Sugar().Errorf("[golog.gin.mongoRecordXml]：%s", err)
	}
}

// MongoDelete 删除
func (c *GinClient) MongoDelete(ctx context.Context, hour int64) (*mongo.DeleteResult, error) {
	filter := bson.D{{"log_time", bson.D{{"$lt", primitive.NewDateTimeFromTime(gotime.Current().BeforeHour(hour).Time)}}}}
	return c.mongoClient.Db.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).DeleteMany(ctx, filter)
}
