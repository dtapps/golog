package golog

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 结构体模型
type ginMongoLog struct {
	LogID          primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                        //【记录】编号
	LogTime        primitive.DateTime `json:"log_time,omitempty" bson:"log_time"`                           //【记录】时间
	TraceID        string             `json:"trace_id,omitempty" bson:"trace_id,omitempty"`                 //【记录】跟踪编号
	RequestTime    string             `json:"request_time,omitempty" bson:"request_time,omitempty"`         //【请求】时间
	RequestUri     string             `json:"request_uri,omitempty" bson:"request_uri,omitempty"`           //【请求】链接 域名+路径+参数
	RequestURL     string             `json:"request_url,omitempty" bson:"request_url,omitempty"`           //【请求】链接 域名+路径
	RequestApi     string             `json:"request_api,omitempty" bson:"request_api,omitempty"`           //【请求】接口
	RequestMethod  string             `json:"request_method,omitempty" bson:"request_method,omitempty"`     //【请求】方式
	RequestProto   string             `json:"request_proto,omitempty" bson:"request_proto,omitempty"`       //【请求】协议
	RequestBody    interface{}        `json:"request_body,omitempty" bson:"request_body,omitempty"`         //【请求】参数
	RequestIP      string             `json:"request_ip,omitempty" bson:"request_ip,omitempty"`             //【请求】客户端IP
	RequestHeader  string             `json:"request_header,omitempty" bson:"request_header,omitempty"`     //【请求】头部
	ResponseTime   time.Time          `json:"response_time,omitempty" bson:"response_time,omitempty"`       //【返回】时间
	ResponseCode   int                `json:"response_code,omitempty" bson:"response_code,omitempty"`       //【返回】状态码
	ResponseData   interface{}        `json:"response_data,omitempty" bson:"response_data,omitempty"`       //【返回】数据
	CostTime       int64              `json:"cost_time,omitempty" bson:"cost_time,omitempty"`               //【系统】花费时间
	SystemHostName string             `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"` //【系统】主机名
	SystemInsideIP string             `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"` //【系统】内网IP
	SystemOs       string             `json:"system_os,omitempty" bson:"system_os,omitempty"`               //【系统】类型
	SystemArch     string             `json:"system_arch,omitempty" bson:"system_arch,omitempty"`           //【系统】架构
	SystemUpTime   uint64             `json:"system_up_time,omitempty" bson:"system_up_time,omitempty"`     //【系统】运行时间
	SystemBootTime uint64             `json:"system_boot_time,omitempty" bson:"system_boot_time,omitempty"` //【系统】开机时间
	GoVersion      string             `json:"go_version,omitempty" bson:"go_version,omitempty"`             //【程序】Go版本
	SdkVersion     string             `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`           //【程序】Sdk版本
	SystemVersion  string             `json:"system_version,omitempty" bson:"system_version,omitempty"`     //【程序】System版本
	CpuCores       int                `json:"cpu_cores,omitempty" bson:"cpu_cores,omitempty"`               //【CPU】核数
	CpuModelName   string             `json:"cpu_model_name,omitempty" bson:"cpu_model_name,omitempty"`     //【CPU】型号名称
	CpuMhz         float64            `json:"cpu_mhz,omitempty" bson:"cpu_mhz,omitempty"`                   //【CPU】兆赫
}

// 创建时间序列集合
func (gm *GinMongo) mongoCreateCollection(ctx context.Context) {
	if gm.mongoConfig.stats == false {
		return
	}

	err := gm.mongoClient.Database(gm.mongoConfig.databaseName).
		CreateCollection(ctx,
			gm.mongoConfig.collectionName,
			options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField("log_time")))
	if err != nil {
		if gm.slog.status {
			gm.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("创建时间序列集合：%s", err))
		}
	}
}

// 创建索引
func (gm *GinMongo) mongoCreateIndexes(ctx context.Context) {
	if gm.mongoConfig.stats == false {
		return
	}

	_, err := gm.mongoClient.Database(gm.mongoConfig.databaseName).
		Collection(gm.mongoConfig.collectionName).
		Indexes().
		CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{{
					Key:   "log_time",
					Value: -1,
				}},
			}})
	if err != nil {
		if gm.slog.status {
			gm.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("创建索引：%s", err))
		}
	}
}
