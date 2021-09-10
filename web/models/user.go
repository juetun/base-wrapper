// Package models
package models

import (
	"github.com/juetun/base-wrapper/lib/base"
)

type User struct {
	base.Model
	UserHid string `json:"user_hid"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
}

func (r *User) GetUserHid() (userHid string, err error) {
	userHid = r.UserHid
	return
}

func (User) TableName() string {
	return "user_main"
}
