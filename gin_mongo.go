package golog

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
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
	LogId           primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                          //【记录】编号
	TraceId         string             `json:"trace_id,omitempty" bson:"trace_id,omitempty"`                   //【记录】跟踪编号
	RequestTime     dorm.BsonTime      `json:"request_time,omitempty" bson:"request_time,omitempty"`           //【请求】时间
	RequestUri      string             `json:"request_uri,omitempty" bson:"request_uri,omitempty"`             //【请求】请求链接 域名+路径+参数
	RequestUrl      string             `json:"request_url,omitempty" bson:"request_url,omitempty"`             //【请求】请求链接 域名+路径
	RequestApi      string             `json:"request_api,omitempty" bson:"request_api,omitempty"`             //【请求】请求接口 路径
	RequestMethod   string             `json:"request_method,omitempty" bson:"request_method,omitempty"`       //【请求】请求方式
	RequestProto    string             `json:"request_proto,omitempty" bson:"request_proto,omitempty"`         //【请求】请求协议
	RequestUa       string             `json:"request_ua,omitempty" bson:"request_ua,omitempty"`               //【请求】请求UA
	RequestReferer  string             `json:"request_referer,omitempty" bson:"request_referer,omitempty"`     //【请求】请求referer
	RequestBody     interface{}        `json:"request_body,omitempty" bson:"request_body,omitempty"`           //【请求】请求主体
	RequestUrlQuery interface{}        `json:"request_url_query,omitempty" bson:"request_url_query,omitempty"` //【请求】请求URL参数
	RequestIp       struct {
		Ip        string `json:"ip,omitempty" bson:"ip,omitempty"`               //【请求】请求客户端Ip
		Continent string `json:"continent,omitempty" bson:"continent,omitempty"` //【请求】请求客户端大陆
		Country   string `json:"country,omitempty" bson:"country,omitempty"`     //【请求】请求客户端国家
		Province  string `json:"province,omitempty" bson:"province,omitempty"`   //【请求】请求客户端省份
		City      string `json:"city,omitempty" bson:"city,omitempty"`           //【请求】请求客户端城市
		Isp       string `json:"isp,omitempty" bson:"isp,omitempty"`             //【请求】请求客户端运营商
	} `json:"request_ip,omitempty" bson:"request_ip,omitempty"` //【请求】请求客户端信息
	RequestIpLocation interface{}   `json:"request_ip_location,omitempty" bson:"request_ip_location,omitempty"` //【请求】请求客户端位置
	RequestHeader     interface{}   `json:"request_header,omitempty" bson:"request_header,omitempty"`           //【请求】请求头
	ResponseTime      dorm.BsonTime `json:"response_time,omitempty" bson:"response_time,omitempty"`             //【返回】时间
	ResponseCode      int           `json:"response_code,omitempty" bson:"response_code,omitempty"`             //【返回】状态码
	ResponseMsg       string        `json:"response_msg,omitempty" bson:"response_msg,omitempty"`               //【返回】描述
	ResponseData      interface{}   `json:"response_data,omitempty" bson:"response_data,omitempty"`             //【返回】数据
	CostTime          int64         `json:"cost_time,omitempty" bson:"cost_time,omitempty"`                     //【系统】花费时间
	System            struct {
		HostName  string `json:"host_name" bson:"host_name"`   //【系统】主机名
		InsideIp  string `json:"inside_ip" bson:"inside_ip"`   //【系统】内网ip
		OutsideIp string `json:"outside_ip" bson:"outside_ip"` //【系统】外网ip
		Os        string `json:"os" bson:"os"`                 //【系统】系统类型
		Arch      string `json:"arch" bson:"arch"`             //【系统】系统架构
	} `json:"system" bson:"system"` //【系统】信息
	Version struct {
		Go  string `json:"go" bson:"go"`   //【程序】Go版本
		Sdk string `json:"sdk" bson:"sdk"` //【程序】Sdk版本
	} `json:"version" bson:"version"` //【程序】版本信息
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
	_, err := c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).CreateManyIndexes(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{
				Key:   "trace_id",
				Value: -1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_time",
				Value: -1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "request_method",
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
				Key:   "request_ip_location",
				Value: "2dsphere",
			}},
		},
	})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建索引：%s", err)
	}
}

// MongoDelete 删除
func (c *GinClient) MongoDelete(ctx context.Context, hour int64) (*mongo.DeleteResult, error) {
	filter := bson.D{{"request_time", bson.D{{"$lt", dorm.NewBsonTimeFromTime(gotime.Current().BeforeHour(hour).Time)}}}}
	return c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).DeleteMany(ctx, filter)
}

// 记录日志
func (c *GinClient) mongoRecord(ctx context.Context, data ginMongoLog) {

	data.LogId = primitive.NewObjectID()             //【记录】编号
	data.System.HostName = c.config.systemHostName   //【系统】主机名
	data.System.InsideIp = c.config.systemInsideIp   //【系统】内网ip
	data.System.OutsideIp = c.config.systemOutsideIp //【系统】外网ip
	data.System.Os = c.config.systemOs               //【系统】系统类型
	data.System.Arch = c.config.systemArch           //【系统】系统架构
	data.Version.Go = c.config.goVersion             //【程序】Go版本
	data.Version.Sdk = c.config.sdkVersion           //【程序】Sdk版本

	_, err := c.mongoClient.Database(c.mongoConfig.databaseName).Collection(c.mongoConfig.collectionName).InsertOne(ctx, data)
	if err != nil {
		c.zapLog.WithTraceIdStr(data.TraceId).Sugar().Errorf("保存框架日志错误：%s", err)
		c.zapLog.WithTraceIdStr(data.TraceId).Sugar().Errorf("保存框架日志数据：%+v", data)
	}
}

func (c *GinClient) mongoRecordJson(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, ipInfo goip.AnalyseResult) {

	var ctx = gotrace_id.SetGinTraceIdContext(context.Background(), ginCtx)

	data := ginMongoLog{
		TraceId:         traceId,                                                      //【记录】跟踪编号
		RequestTime:     dorm.NewBsonTimeFromTime(requestTime),                        //【请求】时间
		RequestUrl:      ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:      gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:   ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:    ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:       ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:  ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery: ginCtx.Request.URL.Query(),                                   //【请求】请求URL参数
		RequestHeader:   ginCtx.Request.Header,                                        //【请求】请求头
		ResponseTime:    dorm.NewBsonTimeCurrent(),                                    //【返回】时间
		ResponseCode:    responseCode,                                                 //【返回】状态码
		ResponseData:    dorm.JsonDecodeNoError([]byte(responseBody)),                 //【返回】数据
		CostTime:        endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBody = dorm.JsonDecodeNoError(requestBody) //【请求】请求主体
	}

	//【请求】请求客户端信息
	data.RequestIp.Ip = ipInfo.Ip
	data.RequestIp.Continent = ipInfo.Continent
	data.RequestIp.Country = ipInfo.Country
	data.RequestIp.Province = ipInfo.Province
	data.RequestIp.City = ipInfo.City
	data.RequestIp.City = ipInfo.Isp
	if ipInfo.LocationLatitude != 0 && ipInfo.LocationLongitude != 0 {
		data.RequestIpLocation = ginMongoLogRequestIpLocationLocation{
			Type:        "Point",
			Coordinates: []float64{ipInfo.LocationLongitude, ipInfo.LocationLatitude},
		}
	}

	c.mongoRecord(ctx, data)
}

func (c *GinClient) mongoRecordXml(ginCtx *gin.Context, traceId string, requestTime time.Time, requestBody []byte, responseCode int, responseBody string, startTime, endTime int64, ipInfo goip.AnalyseResult) {

	var ctx = gotrace_id.SetGinTraceIdContext(context.Background(), ginCtx)

	data := ginMongoLog{
		TraceId:         traceId,                                                      //【记录】跟踪编号
		RequestTime:     dorm.NewBsonTimeFromTime(requestTime),                        //【请求】时间
		RequestUrl:      ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:      gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod:   ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:    ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestUa:       ginCtx.Request.UserAgent(),                                   //【请求】请求UA
		RequestReferer:  ginCtx.Request.Referer(),                                     //【请求】请求referer
		RequestUrlQuery: ginCtx.Request.URL.Query(),                                   //【请求】请求URL参数
		RequestHeader:   ginCtx.Request.Header,                                        //【请求】请求头
		ResponseTime:    dorm.NewBsonTimeCurrent(),                                    //【返回】时间
		ResponseCode:    responseCode,                                                 //【返回】状态码
		ResponseData:    dorm.JsonDecodeNoError([]byte(responseBody)),                 //【返回】数据
		CostTime:        endTime - startTime,                                          //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}

	if len(requestBody) > 0 {
		data.RequestBody = dorm.XmlDecodeNoError(requestBody) //【请求】请求主体
	}

	//【请求】请求客户端信息
	data.RequestIp.Ip = ipInfo.Ip
	data.RequestIp.Continent = ipInfo.Continent
	data.RequestIp.Country = ipInfo.Country
	data.RequestIp.Province = ipInfo.Province
	data.RequestIp.City = ipInfo.City
	data.RequestIp.City = ipInfo.Isp
	if ipInfo.LocationLatitude != 0 && ipInfo.LocationLongitude != 0 {
		data.RequestIpLocation = ginMongoLogRequestIpLocationLocation{
			Type:        "Point",
			Coordinates: []float64{ipInfo.LocationLongitude, ipInfo.LocationLatitude},
		}
	}

	c.mongoRecord(ctx, data)
}
