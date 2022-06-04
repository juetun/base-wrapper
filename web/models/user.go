// Package models
package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type User struct {
	base.Model
	UserHid int64  `json:"user_hid"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
}

func (r *User) GetUserHid() (userHid string, err error) {
	userHid = fmt.Sprintf("%d", r.UserHid)
	return
}

func (User) TableName() string {
	return "user_main"
}
