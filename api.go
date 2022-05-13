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
	LogId                 uint           `gorm:"primaryKey" json:"log_id"`                   //【记录】编号
	RequestTime           TimeString     `gorm:"index" json:"request_time"`                  //【请求】时间
	RequestUri            string         `gorm:"type:text" json:"request_uri"`               //【请求】链接
	RequestUrl            string         `gorm:"type:text" json:"request_url"`               //【请求】链接
	RequestApi            string         `gorm:"type:text;index" json:"request_api"`         //【请求】接口
	RequestMethod         string         `gorm:"type:text;index" json:"request_method"`      //【请求】方式
	RequestParams         datatypes.JSON `gorm:"type:jsonb" json:"request_params"`           //【请求】参数
	RequestHeader         datatypes.JSON `gorm:"type:jsonb" json:"request_header"`           //【请求】头部
	ResponseHeader        datatypes.JSON `gorm:"type:jsonb" json:"response_header"`          //【返回】头部
	ResponseStatusCode    int            `gorm:"type:bigint" json:"response_status_code"`    //【返回】状态码
	ResponseBody          datatypes.JSON `gorm:"type:jsonb" json:"response_body"`            //【返回】内容
	ResponseContentLength int64          `gorm:"type:bigint" json:"response_content_length"` //【返回】大小
	ResponseTime          TimeString     `gorm:"index" json:"response_time"`                 //【返回】时间
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

// Query 查询
func (a *api) Query() *gorm.DB {
	return a.db.Table(a.tableName)
}
