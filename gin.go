package golog

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// 框架定义
type gin struct {
	db        *gorm.DB // pgsql数据库
	tableName string   // 日志表名
}

// GinPostgresqlLog 结构体
type GinPostgresqlLog struct {
	LogId             uint           `gorm:"primaryKey"`      //【记录】编号
	TraceId           string         `gorm:"type:text"`       //【系统】链编号
	RequestTime       time.Time      `gorm:"index"`           //【请求】时间
	RequestUri        string         `gorm:"type:text"`       //【请求】请求链接 域名+路径+参数
	RequestUrl        string         `gorm:"type:text"`       //【请求】请求链接 域名+路径
	RequestApi        string         `gorm:"type:text;index"` //【请求】请求接口 路径
	RequestMethod     string         `gorm:"type:text;index"` //【请求】请求方式
	RequestProto      string         `gorm:"type:text"`       //【请求】请求协议
	RequestUa         string         `gorm:"type:text"`       //【请求】请求UA
	RequestReferer    string         `gorm:"type:text"`       //【请求】请求referer
	RequestBody       datatypes.JSON `gorm:"type:jsonb"`      //【请求】请求主体
	RequestUrlQuery   datatypes.JSON `gorm:"type:jsonb"`      //【请求】请求URL参数
	RequestIp         string         `gorm:"type:text"`       //【请求】请求客户端Ip
	RequestIpCountry  string         `gorm:"type:text"`       //【请求】请求客户端城市
	RequestIpRegion   string         `gorm:"type:text"`       //【请求】请求客户端区域
	RequestIpProvince string         `gorm:"type:text"`       //【请求】请求客户端省份
	RequestIpCity     string         `gorm:"type:text"`       //【请求】请求客户端城市
	RequestIpIsp      string         `gorm:"type:text"`       //【请求】请求客户端运营商
	RequestHeader     datatypes.JSON `gorm:"type:jsonb"`      //【请求】请求头
	ResponseTime      time.Time      `gorm:"index"`           //【返回】时间
	ResponseCode      int            `gorm:"type:bigint"`     //【返回】状态码
	ResponseMsg       string         `gorm:"type:text"`       //【返回】描述
	ResponseData      datatypes.JSON `gorm:"type:jsonb"`      //【返回】数据
	CostTime          int64          `gorm:"type:bigint"`     //【系统】花费时间
}

// AutoMigrate 自动迁移
func (g *gin) AutoMigrate() {
	err := g.db.Table(g.tableName).AutoMigrate(&GinPostgresqlLog{})
	if err != nil {
		panic("创建表失败：" + err.Error())
	}
}

// Record 记录日志
func (g *gin) Record(content GinPostgresqlLog) int64 {
	return g.db.Table(g.tableName).Create(&content).RowsAffected
}
