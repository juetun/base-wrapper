package base

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TimeNormalDate struct {
	time.Time
}

func (t TimeNormalDate) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02"`)
	return []byte(tune), nil
}

// IsZero 判断当前时间是否为空
func (t *TimeNormalDate) IsZero() (res bool) {
	if t.Time.IsZero() {
		res = true
	}
	return
}

func (t *TimeNormalDate) Format(layout string) (res string) {
	res = t.Time.Format(layout)
	return
}

func (t *TimeNormalDate) UnmarshalJSON(b []byte) (err error) {
	b = bytes.Trim(b, "\"") // 此除需要去掉传入的数据的两端的 ""
	if b == nil {
		return
	}
	v := string(b)
	if t == nil {
		t = &TimeNormalDate{}
	}
	t.Time, err = time.ParseInLocation("2006-01-02", v, time.Local)
	return
}
func (t TimeNormalDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *TimeNormalDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = TimeNormalDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *TimeNormalDate) BeforeCreate(scope *gorm.DB) error {
	timeNow := time.Now()
	scope.Set("create_time", timeNow)
	scope.Set("update_time", timeNow)
	return nil
}

func (t *TimeNormalDate) BeforeUpdate(scope *gorm.DB) error {
	scope.Set("update_time", time.Now())
	return nil
}
