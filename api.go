package golog

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"log"
)

// 接口定义
type api struct {
	db        *gorm.DB // pgsql数据库
	tableName string   // 日志表名
	insideIp  string   // 内网ip
	hostname  string   // 主机名
	goVersion string   // go版本
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
	SystemHostName        string         `gorm:"type:text" json:"system_host_name"`          //【系统】主机名
	SystemInsideIp        string         `gorm:"type:text" json:"system_inside_ip"`          //【系统】内网ip
	GoVersion             string         `gorm:"type:text" json:"go_version"`                //【程序】Go版本
}

// AutoMigrate 自动迁移
func (p *api) AutoMigrate() {
	err := p.db.Table(p.tableName).AutoMigrate(&ApiPostgresqlLog{})
	if err != nil {
		panic("创建表失败：" + err.Error())
	}
}

// Record 记录日志
func (p *api) Record(content ApiPostgresqlLog) *gorm.DB {
	content.SystemHostName = p.hostname
	if content.SystemInsideIp == "" {
		content.SystemInsideIp = p.insideIp
	}
	content.GoVersion = p.goVersion
	resp := p.db.Table(p.tableName).Create(&content)
	if resp.RowsAffected == 0 {
		log.Println("api：", resp.Error)
	}
	return resp
}

// Query 查询
func (p *api) Query() *gorm.DB {
	return p.db.Table(p.tableName)
}
