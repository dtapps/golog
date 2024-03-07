package golog

import (
	"context"
	"log"
)

// gormRecord 记录日志
func (hg *HertzGorm) gormRecord(ctx context.Context, data hertzGormLog) {
	if hg.gormConfig.stats == false {
		return
	}

	err := hg.gormClient.WithContext(ctx).
		Table(hg.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		log.Printf("记录接口日志错误：%s", err)
		log.Printf("记录接口日志数据：%+v", data)
	}
}
