// Package models
package models

import (
	"fmt"

	"github.com/juetun/base-wrapper/lib/base"
)

type User struct {
	base.Model
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (r *User) GetUserHid() (userHid string, err error) {
	userHid = fmt.Sprintf("%v", r.Id)
	return
}

func (User) TableName() string {
	return "user_main"
}
