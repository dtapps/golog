package golog

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// 结构体模型
type hertzGormLog struct {
	RequestID          string    `gorm:"index;comment:请求编号" json:"request_id,omitempty"`     // 请求编号
	RequestTime        time.Time `gorm:"index;comment:请求时间" json:"request_time,omitempty"`   // 请求时间
	RequestHost        string    `gorm:"comment:请求主机" json:"request_host,omitempty"`         // 请求主机
	RequestPath        string    `gorm:"index;comment:请求地址" json:"request_path,omitempty"`   // 请求地址
	RequestQuery       string    `gorm:"comment:请求参数" json:"request_query,omitempty"`        // 请求参数
	RequestMethod      string    `gorm:"index;comment:请求方式" json:"request_method,omitempty"` // 请求方式
	RequestScheme      string    `gorm:"comment:请求协议" json:"request_scheme,omitempty"`       // 请求协议
	RequestContentType string    `gorm:"comment:请求类型" json:"request_content_type,omitempty"` // 请求类型
	RequestBody        string    `gorm:"comment:请求内容" json:"request_body,omitempty"`         // 请求内容
	RequestClientIP    string    `gorm:"comment:请求IP" json:"request_client_ip,omitempty"`    // 请求IP
	RequestUserAgent   string    `gorm:"comment:请求UA" json:"request_user_agent,omitempty"`   // 请求UA
	RequestHeader      string    `gorm:"comment:请求头" json:"request_header,omitempty"`        // 请求头
	RequestCostTime    float64   `gorm:"comment:请求消耗时长" json:"request_cost_time,omitempty"`  // 请求消耗时长
	ResponseTime       time.Time `gorm:"comment:响应时间" json:"response_time,omitempty"`        // 响应时间
	ResponseHeader     string    `gorm:"comment:响应头" json:"response_header,omitempty"`       // 响应头
	ResponseStatusCode int       `gorm:"comment:响应状态" json:"response_status_code,omitempty"` // 响应状态
	ResponseBody       string    `gorm:"comment:响应内容" json:"response_data,omitempty"`        // 响应内容
	GoVersion          string    `gorm:"comment:Go版本" json:"go_version,omitempty"`           // Go版本
	SdkVersion         string    `gorm:"comment:Sdk版本" json:"sdk_version,omitempty"`         // Sdk版本
	SystemInfo         string    `gorm:"comment:系统信息" json:"system_info,omitempty"`          // 系统信息
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
		slog.Error(fmt.Sprintf("创建模型：%s", err))
	}
}
