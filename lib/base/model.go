/**
* @Author:changjiang
* @Description:
* @File:model
* @Version: 1.0.0
* @Date 2020/3/29 2:25 下午
 */
package base

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Model struct {
	Id        int       `gorm:"primary_key" json:"id"`
	CreatedAt TimeNormal  `json:"created_at"`
	UpdatedAt TimeNormal  `json:"updated_at"`
	DeletedAt *TimeNormal `sql:"index" json:"-"`
}


type TimeNormal struct {
	time.Time
}

func (t TimeNormal) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
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

func (t *TimeNormal) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("create_time", time.Now())
	scope.SetColumn("update_time", time.Now())
	return nil
}

func (t *TimeNormal) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("update_time", time.Now())
	return nil
}
