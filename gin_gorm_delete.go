package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
)

// GormDeleteData 删除N天前数据
func (gg *GinGorm) GormDeleteData(ctx context.Context, day int) error {
	return gg.GormDeleteDataCustom(ctx, gg.gormConfig.tableName, day)
}

// GormDeleteDataCustom 删除N天前数据
func (gg *GinGorm) GormDeleteDataCustom(ctx context.Context, tableName string, day int) error {
	if gg.gormConfig.stats == false {
		return nil
	}

	if tableName == "" {
		return errors.New("没有设置表名")
	}
	return gg.gormClient.Table(tableName).Where("request_time < ?", gotime.Current().BeforeDay(day).Format()).Delete(&ginPostgresqlLog{}).Error
}
