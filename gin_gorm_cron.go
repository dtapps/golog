package golog

import "context"
import "github.com/robfig/cron/v3"

// StartCronClean 定时清理框架日志
func (gg *GinGorm) StartCronClean(ctx context.Context, cr *cron.Cron, cp string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = gg.GormDeleteData(ctx, hour)
	})
}

// StartCustomCronClean 定时清理框架日志
func (gg *GinGorm) StartCustomCronClean(ctx context.Context, cr *cron.Cron, cp string, tableName string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = gg.GormDeleteDataCustom(ctx, tableName, hour)
	})
}
