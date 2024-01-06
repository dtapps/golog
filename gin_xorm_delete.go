package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
)

// XormDeleteData 删除N天前数据
func (gg *GinXorm) XormDeleteData(ctx context.Context, day int) error {
	return gg.XormDeleteDataCustom(ctx, gg.xormConfig.tableName, day)
}

// XormDeleteDataCustom 删除N天前数据r
func (gg *GinXorm) XormDeleteDataCustom(ctx context.Context, tableName string, day int) error {
	if gg.xormConfig.stats == false {
		return nil
	}

	if tableName == "" {
		return errors.New("没有设置表名")
	}
	_, err := gg.xormClient.Table(tableName).Where("request_time < ?", gotime.Current().BeforeDay(day).Format()).Delete(&ginXormLog{})
	return err
}
