// Package app_obj
/**
* @Author:changjiang
* @Description:
* @File:db
* @Version: 1.0.0
* @Date 2020/3/27 10:39 下午
 */
package app_obj

import (
	"github.com/jinzhu/gorm"
)

var DbMysql = make(map[string]*gorm.DB, 2)
