package golog

import (
	"context"
	"log"
	"time"
)

// 结构体模型
type hertzGormLog struct {
	RequestID          string    `json:"request_id" bson:""`           //【日志】ID
	RequestTime        time.Time `json:"request_time" bson:""`         //【请求】Time
	RequestHost        string    `json:"request_host" bson:""`         //【请求】Host
	RequestPath        string    `json:"request_path" bson:""`         //【请求】Path
	RequestQuery       any       `json:"request_query" bson:""`        //【请求】Query
	RequestMethod      string    `json:"request_method" bson:""`       //【请求】Method
	RequestScheme      string    `json:"request_scheme" bson:""`       //【请求】Scheme
	RequestContentType string    `json:"request_content_type" bson:""` //【请求】Content-Type
	RequestBody        any       `json:"request_body" bson:""`         //【请求】Body
	RequestClientIP    string    `json:"request_client_ip" bson:""`    //【请求】ClientIP
	RequestUserAgent   string    `json:"request_user_agent" bson:""`   //【请求】User-Agent
	RequestHeader      any       `json:"request_header" bson:""`       //【请求】Header
	RequestCostTime    int64     `json:"request_cost_time" bson:""`    //【请求】Cost
	ResponseTime       time.Time `json:"response_time" bson:""`        //【响应】Time
	ResponseHeader     any       `json:"response_header" bson:""`      //【响应】Header
	ResponseStatusCode int       `json:"response_status_code" bson:""` //【响应】StatusCode
	ResponseBody       any       `json:"response_data" bson:""`        //【响应】Body
}

// 创建模型
func (hg *HertzGorm) gormAutoMigrate(ctx context.Context) {
	if hg.gormConfig.stats == false {
		return
	}

	err := hg.gormClient.WithContext(ctx).
		Table(hg.gormConfig.tableName).
		AutoMigrate(&hertzGormLog{})
	if err != nil {
		log.Printf("创建模型：%s\n", err)
	}
}
