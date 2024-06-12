package golog

import (
	"context"
	"fmt"
	"log/slog"
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
		slog.Error(fmt.Sprintf("创建模型：%s", err))
	}
}
