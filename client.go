package golog

import "gorm.io/gorm"

const Version = "1.0.0"

type Client struct {
	Gin gin // 框架日志
	Api api // 接口日志
}

// NewClientGin 创建框架实例化
func NewClientGin(db *gorm.DB, tableName string) *Client {
	if db == nil {
		panic("驱动不正常")
	}
	if tableName == "" {
		panic("表名不能为空")
	}
	client := &Client{Gin: gin{db: db, tableName: tableName}}
	client.Gin.AutoMigrate()
	return client
}

// NewClientApi 创建接口实例化
func NewClientApi(db *gorm.DB, tableName string) *Client {
	if db == nil {
		panic("驱动不正常")
	}
	if tableName == "" {
		panic("表名不能为空")
	}
	client := &Client{Api: api{db: db, tableName: tableName}}
	client.Api.AutoMigrate()
	return client
}
