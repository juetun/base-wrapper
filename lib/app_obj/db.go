/**
* @Author:changjiang
* @Description:
* @File:db
* @Version: 1.0.0
* @Date 2020/3/27 10:39 下午
 */
package app_obj

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var DbMysql = make(map[string]*gorm.DB)

// 获取Redis操作实例
func GetDbClient(nameSpace ...string) *gorm.DB {

	var s string
	switch len := len(nameSpace); len {
	case 0:
		s = "default"
	case 1:
		s = nameSpace[0]
	default:
		panic("nameSpace receive at most one parameter")
	}
	if _, ok := DbMysql[s]; ok {
		return DbMysql[s]
	}
	panic(fmt.Sprintf("the Database connect(%s) is not exist", s))
}
