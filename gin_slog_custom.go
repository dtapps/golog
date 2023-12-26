package golog

import (
	"context"
	"go.dtapp.net/goip"
)

// GinSLogCustom 框架自定义日志
type GinSLogCustom struct {
	ipService *goip.Client // IP服务
	slog      struct {
		status bool  // 状态
		client *SLog // 日志服务
	}
}

// GinSLogCustomFun  框架自定义日志驱动
type GinSLogCustomFun func() *GinSLogCustom

// NewGinSLogCustom 创建框架实例化
func NewGinSLogCustom(ctx context.Context, ipService *goip.Client) (*GinSLogCustom, error) {
	c := &GinSLogCustom{}
	c.ipService = ipService
	return c, nil
}
