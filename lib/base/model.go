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

type (
	Model struct {
		Id        int         `gorm:"primary_key" json:"id"`
		CreatedAt TimeNormal  `json:"created_at"`
		UpdatedAt TimeNormal  `json:"updated_at"`
		DeletedAt *TimeNormal `sql:"index" json:"-"`
	}
	CreateTable interface {
		// TableName 获取表名
		TableName() string
		// GetTableComment 获取表注册名
		GetTableComment() (res string)
	}
	TimeNormal struct {
		time.Time
	}
	//字段选项类型（用于前端拼接select  radio checkbox类型数据）
	ModelItemOption struct {
		Value interface{} `json:"value"`
		Label string      `json:"label"`
	}
	ModelItemOptions []ModelItemOption
)

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

func (r *ModelItemOptions) validateValueType(item *ModelItemOption, typeString string) (err error) {
	switch item.Value.(type) {
	case uint8:
		if typeString != "uint8" {
			err = fmt.Errorf("%s格式不正确", item.Label)
			return
		}
	case int8:
		if typeString != "int8" {
			err = fmt.Errorf("%s格式不正确", item.Label)
			return
		}
	case string:
		if typeString != "string" {
			err = fmt.Errorf("%s格式不正确", item.Label)
			return
		}
	case int:
		if typeString != "int" {
			err = fmt.Errorf("%s格式不正确", item.Label)
			return
		}
	case uint64:
		if typeString != "uint64" {
			err = fmt.Errorf("%s格式不正确", item.Label)
			return
		}
	case int64:
		if typeString != "int64" {
			err = fmt.Errorf("%s格式不正确", item.Label)
			return
		}
	default:
		err = fmt.Errorf("当前不支持你选择的类型(%t)", item.Value)
	}
	return
}

func (r *ModelItemOptions) GetMapAsKeyString() (res map[string]string, err error) {
	res = make(map[string]string, len(*r))
	for _, item := range *r {
		if err = r.validateValueType(&item, "string"); err != nil {
			return
		}
		res[item.Value.(string)] = item.Label
	}

	return
}

func (r *ModelItemOptions) GetMapAsKeyUint8() (res map[uint8]string, err error) {
	res = make(map[uint8]string, len(*r))
	for _, item := range *r {
		if err = r.validateValueType(&item, "uint8"); err != nil {
			return
		}
		res[item.Value.(uint8)] = item.Label
	}

	return
}

func (r *ModelItemOptions) GetMapAsKeyInt8() (res map[int8]string, err error) {
	res = make(map[int8]string, len(*r))
	for _, item := range *r {
		if err = r.validateValueType(&item, "int8"); err != nil {
			return
		}
		res[item.Value.(int8)] = item.Label
	}

	return
}

func (r *ModelItemOptions) GetMapAsKeyInt() (res map[int]string, err error) {
	res = make(map[int]string, len(*r))
	for _, item := range *r {
		if err = r.validateValueType(&item, "int"); err != nil {
			return
		}
		res[item.Value.(int)] = item.Label
	}

	return
}

func (r *ModelItemOptions) GetMapAsKeyInt64() (res map[int64]string, err error) {
	res = make(map[int64]string, len(*r))
	for _, item := range *r {
		if err = r.validateValueType(&item, "int64"); err != nil {
			return
		}
		res[item.Value.(int64)] = item.Label
	}

	return
}

func (r *ModelItemOptions) GetMapAsKeyUint64() (res map[uint64]string, err error) {
	res = make(map[uint64]string, len(*r))
	for _, item := range *r {
		if err = r.validateValueType(&item, "uint64"); err != nil {
			return
		}
		res[item.Value.(uint64)] = item.Label
	}

	return
}
