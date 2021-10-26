// Package base
/**
* @Author:changjiang
* @Description:
* @File:model
* @Version: 1.0.0
* @Date 2020/3/29 2:25 下午
 */
package base

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	Id        int         `gorm:"primary_key" json:"id"`
	CreatedAt TimeNormal  `json:"created_at"`
	UpdatedAt TimeNormal  `json:"updated_at"`
	DeletedAt *TimeNormal `sql:"index" json:"-"`
}

type CreateTable interface {
	// TableName 获取表名
	TableName() string
	// GetTableComment 获取表注册名
	GetTableComment() (res string)
}

type TimeNormal struct {
	time.Time
}

// GetNowTimeNormal 获取当前时间的时间格式
func GetNowTimeNormal() (res TimeNormal) {
	res = TimeNormal{Time: time.Now()}
	return
}
func (t TimeNormal) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

// IsZero 判断当前时间是否为空
func (t *TimeNormal) IsZero() (res bool) {
	if t.Time.IsZero() {
		res = true
	}
	return
}
func (t *TimeNormal) UnmarshalJSON(b []byte) (err error) {
	b = bytes.Trim(b, "\"") // 此除需要去掉传入的数据的两端的 ""
	v := string(b)
	if t == nil {
		t = &TimeNormal{}
	}
	t.Time, err = time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
	return
}
func (t TimeNormal) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *TimeNormal) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = TimeNormal{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *TimeNormal) BeforeCreate(scope *gorm.DB) error {
	timeNow := time.Now()
	scope.Set("create_time", timeNow)
	scope.Set("update_time", timeNow)
	return nil
}

func (t *TimeNormal) BeforeUpdate(scope *gorm.DB) error {
	scope.Set("update_time", time.Now())
	return nil
}
