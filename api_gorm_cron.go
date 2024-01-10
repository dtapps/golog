package golog

import "context"
import "github.com/robfig/cron/v3"

// StartCronClean 定时清理接口日志
func (ag *ApiGorm) StartCronClean(ctx context.Context, cr *cron.Cron, cp string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = ag.GormDeleteData(ctx, hour)
	})
}
