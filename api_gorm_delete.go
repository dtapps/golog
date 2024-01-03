package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
)

// GormDeleteData 删除N天前数据
func (ag *ApiGorm) GormDeleteData(ctx context.Context, day int) error {
	return ag.GormDeleteDataCustom(ctx, ag.gormConfig.tableName, day)
}

// GormDeleteDataCustom 删除N天前数据
func (ag *ApiGorm) GormDeleteDataCustom(ctx context.Context, tableName string, day int) error {
	if ag.gormConfig.stats == false {
		return nil
	}

	if tableName == "" {
		return errors.New("没有设置表名")
	}
	return ag.gormClient.Table(tableName).Where("request_time < ?", gotime.Current().BeforeDay(day).Format()).Delete(&apiPostgresqlLog{}).Error
}
