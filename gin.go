package golog

import (
	"go.dtapp.net/goip"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 框架定义
type gin struct {
	db        *gorm.DB // pgsql数据库
	tableName string   // 日志表名
	outsideIp string   // 外网ip
	insideIp  string   // 内网ip
}

// GinPostgresqlLog 结构体
type GinPostgresqlLog struct {
	LogId             uint           `gorm:"primaryKey" json:"log_id"`              //【记录】编号
	TraceId           string         `gorm:"type:text" json:"trace_id"`             //【系统】链编号
	RequestTime       TimeString     `gorm:"index" json:"request_time"`             //【请求】时间
	RequestUri        string         `gorm:"type:text" json:"request_uri"`          //【请求】请求链接 域名+路径+参数
	RequestUrl        string         `gorm:"type:text" json:"request_url"`          //【请求】请求链接 域名+路径
	RequestApi        string         `gorm:"type:text;index" json:"request_api"`    //【请求】请求接口 路径
	RequestMethod     string         `gorm:"type:text;index" json:"request_method"` //【请求】请求方式
	RequestProto      string         `gorm:"type:text" json:"request_proto"`        //【请求】请求协议
	RequestUa         string         `gorm:"type:text" json:"request_ua"`           //【请求】请求UA
	RequestReferer    string         `gorm:"type:text" json:"request_referer"`      //【请求】请求referer
	RequestBody       datatypes.JSON `gorm:"type:jsonb" json:"request_body"`        //【请求】请求主体
	RequestUrlQuery   datatypes.JSON `gorm:"type:jsonb" json:"request_url_query"`   //【请求】请求URL参数
	RequestIp         string         `gorm:"type:text" json:"request_ip"`           //【请求】请求客户端Ip
	RequestIpCountry  string         `gorm:"type:text" json:"request_ip_country"`   //【请求】请求客户端城市
	RequestIpRegion   string         `gorm:"type:text" json:"request_ip_region"`    //【请求】请求客户端区域
	RequestIpProvince string         `gorm:"type:text" json:"request_ip_province"`  //【请求】请求客户端省份
	RequestIpCity     string         `gorm:"type:text" json:"request_ip_city"`      //【请求】请求客户端城市
	RequestIpIsp      string         `gorm:"type:text" json:"request_ip_isp"`       //【请求】请求客户端运营商
	RequestHeader     datatypes.JSON `gorm:"type:jsonb" json:"request_header"`      //【请求】请求头
	ResponseTime      TimeString     `gorm:"index" json:"response_time"`            //【返回】时间
	ResponseCode      int            `gorm:"type:bigint" json:"response_code"`      //【返回】状态码
	ResponseMsg       string         `gorm:"type:text" json:"response_msg"`         //【返回】描述
	ResponseData      datatypes.JSON `gorm:"type:jsonb" json:"response_data"`       //【返回】数据
	CostTime          int64          `gorm:"type:bigint" json:"cost_time"`          //【系统】花费时间
	SystemOutsideIp   string         `gorm:"type:text" json:"system_outside_ip"`    //【系统】外网ip
	SystemInsideIp    string         `gorm:"type:text" json:"system_inside_ip"`     //【系统】内网ip
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
	if content.SystemOutsideIp == "" {
		content.SystemOutsideIp = g.outsideIp
	}
	if content.SystemInsideIp == "" {
		content.SystemInsideIp = g.insideIp
	}
	return g.db.Table(g.tableName).Create(&content).RowsAffected
}

// Query 查询
func (g *gin) Query() *gorm.DB {
	return g.db.Table(g.tableName)
}

// 配置
func (g *gin) configOutsideIp() {
	go func() {
		g.insideIp = goip.GetOutsideIp()
	}()
}
