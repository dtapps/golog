package golog

import (
	"encoding/json"
	"go.dtapp.net/gotime"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestTest(t *testing.T) {
	dsn := "host=119.29.14.159 user=dbadmin password=98jolg256s.* dbname=logs port=15432 sslmode=disable TimeZone=Asia/Shanghai"
	pgsqlDB, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	client := NewClientApi(pgsqlDB, "test")

	// 插入时间
	client.Api.Record(ApiPostgresqlLog{
		RequestTime:  TimeString{Time: gotime.Current().Time}, //【请求】时间
		ResponseTime: TimeString{Time: gotime.Current().Time}, //【返回】时间
	})

	// 查询数据
	var result ApiPostgresqlLog
	client.Api.Query().Where("log_id = ?", 10).Take(&result)

	t.Log(result)
	t.Logf("result:%v", result)
	t.Logf("result.request_time:%v", result.RequestTime)
	t.Logf("result:%+v", result)
	t.Logf("result.request_time:%+v", result.RequestTime)

	marshal, err := json.Marshal(result)
	t.Logf("Marshal:%s", marshal)
	t.Logf("Marshal:%v", marshal)
	t.Logf("err:%v", err)

	var jsonM ApiPostgresqlLog
	err = json.Unmarshal(marshal, &jsonM)
	t.Logf("Unmarshal:%v", jsonM)
	t.Logf("Unmarshal:%+v", jsonM)
	t.Logf("err:%v", err)
}
