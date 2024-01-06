package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
)

// XormDeleteData 删除N天前数据
func (ag *ApiXorm) XormDeleteData(ctx context.Context, day int) error {
	return ag.XormDeleteDataCustom(ctx, ag.xormConfig.tableName, day)
}

// XormDeleteDataCustom 删除N天前数据
func (ag *ApiXorm) XormDeleteDataCustom(ctx context.Context, tableName string, day int) error {
	if ag.xormConfig.stats == false {
		return nil
	}

	if tableName == "" {
		return errors.New("没有设置表名")
	}
	_, err := ag.xormClient.Table(tableName).Where("request_time < ?", gotime.Current().BeforeDay(day).Format()).Delete(&apiXormLog{})
	return err
}
