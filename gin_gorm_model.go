package golog

import (
	"context"
	"fmt"
	"log/slog"
)

// 创建模型
func (gg *GinGorm) gormAutoMigrate(ctx context.Context) {
	if gg.gormConfig.stats == false {
		return
	}

	err := gg.gormClient.WithContext(ctx).
		Table(gg.gormConfig.tableName).
		AutoMigrate(&GormGinLogModel{})
	if err != nil {
		slog.Error(fmt.Sprintf("创建模型：%s", err))
	}
}
