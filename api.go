package golog

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 接口定义
type api struct {
	db        *gorm.DB // pgsql数据库
	tableName string   // 日志表名
}

// ApiPostgresqlLog 结构体
type ApiPostgresqlLog struct {
	LogId                 uint           `gorm:"primaryKey"`      //【记录】编号
	RequestTime           TimeString     `gorm:"index"`           //【请求】时间
	RequestUri            string         `gorm:"type:text"`       //【请求】链接
	RequestUrl            string         `gorm:"type:text"`       //【请求】链接
	RequestApi            string         `gorm:"type:text;index"` //【请求】接口
	RequestMethod         string         `gorm:"type:text;index"` //【请求】方式
	RequestParams         datatypes.JSON `gorm:"type:jsonb"`      //【请求】参数
	RequestHeader         datatypes.JSON `gorm:"type:jsonb"`      //【请求】头部
	ResponseHeader        datatypes.JSON `gorm:"type:jsonb"`      //【返回】头部
	ResponseStatusCode    int            `gorm:"type:bigint"`     //【返回】状态码
	ResponseBody          datatypes.JSON `gorm:"type:jsonb"`      //【返回】内容
	ResponseContentLength int64          `gorm:"type:bigint"`     //【返回】大小
	ResponseTime          TimeString     `gorm:"index"`           //【返回】时间
}

// AutoMigrate 自动迁移
func (a *api) AutoMigrate() {
	err := a.db.Table(a.tableName).AutoMigrate(&ApiPostgresqlLog{})
	if err != nil {
		panic("创建表失败：" + err.Error())
	}
}

// Record 记录日志
func (a *api) Record(content ApiPostgresqlLog) int64 {
	return a.db.Table(a.tableName).Create(&content).RowsAffected
}
