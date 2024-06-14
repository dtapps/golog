package golog

import (
	"context"
	"go.opentelemetry.io/otel/codes"
)

// 创建模型
func (ag *ApiGorm) gormAutoMigrate(ctx context.Context) {
	if ag.gormConfig.stats == false {
		return
	}

	err := ag.gormClient.WithContext(ctx).
		Table(ag.gormConfig.tableName).
		AutoMigrate(&GormApiLogModel{})
	if err != nil {
		ag.TraceRecordError(err)
		ag.TraceSetStatus(codes.Error, err.Error())
	}
}
