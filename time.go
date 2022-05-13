package golog

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"go.dtapp.net/gotime"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type TimeString time.Time

// GormDataType gorm通用数据类型
func (t TimeString) GormDataType() string {
	return "string"
}

func (t TimeString) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// 使用 field.Tag、field.TagSettings 获取字段的 tag
	// 查看 https://github.com/go-gorm/gorm/blob/master/schema/field.go 获取全部的选项

	// 根据不同的数据库驱动返回不同的数据类型
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "string"
	case "postgres":
		return "string"
	}
	return ""
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Time
func (t *TimeString) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("无法解析:", value))
	}
	t1 := gotime.SetCurrentParse(str).Time
	*t = TimeString(t1)
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 string value
func (t TimeString) Value() (driver.Value, error) {
	return gotime.SetCurrent(time.Time(t)).Format(), nil
}
