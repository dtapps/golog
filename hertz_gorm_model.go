package golog

import (
	"context"
	"fmt"
	"log/slog"
)

// 创建模型
func (hg *HertzGorm) gormAutoMigrate(ctx context.Context) {
	if hg.gormConfig.stats == false {
		return
	}

	err := hg.gormClient.WithContext(ctx).
		Table(hg.gormConfig.tableName).
		AutoMigrate(&GormHertzLogModel{})
	if err != nil {
		slog.Error(fmt.Sprintf("创建模型：%s", err))
	}
}
