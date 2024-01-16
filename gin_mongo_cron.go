package golog

import "context"
import "github.com/robfig/cron/v3"

// StartCronClean 定时清理框架日志
func (gm *GinMongo) StartCronClean(ctx context.Context, cr *cron.Cron, cp string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = gm.MongoDelete(ctx, hour)
	})
}

// StartCustomCronClean 定时清理框架日志
func (gm *GinMongo) StartCustomCronClean(ctx context.Context, cr *cron.Cron, cp string, databaseName string, collectionName string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = gm.MongoDeleteDataCustom(ctx, databaseName, collectionName, hour)
	})
}
