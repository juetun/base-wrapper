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

const defaultNameSpace = "default"

type GetDbClientDataCallBack func(db *gorm.DB) (err error)

// 获取数据库连接 参数结构体
type GetDbClientData struct {
	DbNameSpace string                  `json:"db_name_space"`
	CallBack    GetDbClientDataCallBack // 获取数据库回调信息
}

// 获取Redis操作实例
func GetDbClient(params ...*GetDbClientData) (db *gorm.DB) {
	l := len(params)

	var arg *GetDbClientData
	if l > 1 {
		panic("arg is more than one parameters")
	} else if l == 1 {
		arg = params[0]
	} else {
		arg = &GetDbClientData{DbNameSpace: defaultNameSpace}
	}

	if arg.DbNameSpace == "" {
		arg.DbNameSpace = defaultNameSpace
	}
	var ok bool
	if db, ok = DbMysql[arg.DbNameSpace]; ok {
		arg.CallBack(db)
		return
	}
	panic(fmt.Sprintf("the Database connect(%s) is not exist", arg.DbNameSpace))
}
