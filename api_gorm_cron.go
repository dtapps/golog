package golog

import (
	"context"
	"github.com/robfig/cron/v3"
)

// StartCronClean 定时清理接口日志
func (ag *ApiGorm) StartCronClean(ctx context.Context, cr *cron.Cron, cp string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = ag.GormDeleteData(ctx, hour)
	})
}

// StartCustomCronClean 定时清理接口日志
func (ag *ApiGorm) StartCustomCronClean(ctx context.Context, cr *cron.Cron, cp string, tableName string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = ag.GormDeleteDataCustom(ctx, tableName, hour)
	})
}
