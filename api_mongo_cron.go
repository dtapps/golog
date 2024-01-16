package golog

import (
	"context"
	"github.com/robfig/cron/v3"
)

// StartCronClean 定时清理接口日志
func (am *ApiMongo) StartCronClean(ctx context.Context, cr *cron.Cron, cp string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = am.MongoDelete(ctx, hour)
	})
}

// StartCustomCronClean 定时清理接口日志
func (am *ApiMongo) StartCustomCronClean(ctx context.Context, cr *cron.Cron, cp string, databaseName string, collectionName string, hour int64) (cron.EntryID, error) {
	return cr.AddFunc(cp, func() {
		_ = am.MongoDeleteDataCustom(ctx, databaseName, collectionName, hour)
	})
}
