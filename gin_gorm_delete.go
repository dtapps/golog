package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
)

// DeleteData 删除N天前数据
func (gg *GinGorm) DeleteData(ctx context.Context, day int) error {
	if gg.gormConfig.tableName == "" {
		return errors.New("没有设置表名")
	}
	return gg.gormClient.GetDb().Table(gg.gormConfig.tableName).Where("request_time >= ?", gotime.Current().BeforeDay(day).Time).Delete(&ginPostgresqlLog{}).Error
}

// DeleteDataCustom 删除N天前数据
func (gg *GinGorm) DeleteDataCustom(ctx context.Context, tableName string, day int) error {
	if tableName == "" {
		return errors.New("没有设置表名")
	}
	return gg.gormClient.GetDb().Table(tableName).Where("request_time >= ?", gotime.Current().BeforeDay(day).Time).Delete(&ginPostgresqlLog{}).Error
}
