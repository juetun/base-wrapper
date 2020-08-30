/**
* @Author:changjiang
* @Description:
* @File:user
* @Version: 1.0.0
* @Date 2020/8/18 6:41 下午
 */
package models

import (
	"github.com/juetun/base-wrapper/lib/base"
)

type User struct {
	base.Model
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (User) TableName() string {
	return "user"
}
