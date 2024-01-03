package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
)

// DeleteData 删除N天前数据
func (ag *ApiGorm) DeleteData(ctx context.Context, day int) error {
	if ag.gormConfig.tableName == "" {
		return errors.New("没有设置表名")
	}
	return ag.gormClient.GetDb().Table(ag.gormConfig.tableName).Where("request_time >= ?", gotime.Current().BeforeDay(day).Time).Delete(&apiPostgresqlLog{}).Error
}

// DeleteDataCustom 删除N天前数据
func (ag *ApiGorm) DeleteDataCustom(ctx context.Context, tableName string, day int) error {
	if tableName == "" {
		return errors.New("没有设置表名")
	}
	return ag.gormClient.GetDb().Table(tableName).Where("request_time >= ?", gotime.Current().BeforeDay(day).Time).Delete(&apiPostgresqlLog{}).Error
}
