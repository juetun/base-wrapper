/**
* @Author:changjiang
* @Description:
* @File:db
* @Version: 1.0.0
* @Date 2021/1/30 4:42 下午
 */
package base

import (
	"fmt"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"gorm.io/gorm"
)

type GetDbClientDataCallBack func(db *gorm.DB) (err error)

const defaultNameSpace = "default"

// 获取数据库连接 参数结构体
type GetDbClientData struct {
	Context     *Context
	DbNameSpace string                  `json:"db_name_space"`
	CallBack    GetDbClientDataCallBack // 获取数据库回调信息
}

func (r *GetDbClientData) DefaultGetDbClientDataCallBack(db *gorm.DB) (err error) {
	var s string
	if nil != r.Context.GinContext {
		if tp, ok := r.Context.GinContext.Get(app_obj.TraceId); ok {
			s = fmt.Sprintf("%v", tp)
		}
	}
	db.InstanceSet(app_obj.TraceId, s)
	// db.InstantSet(app_obj.TraceId, s)
	return

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
	if db, ok = app_obj.DbMysql[arg.DbNameSpace]; ok {
		if arg.CallBack == nil {
			arg.DefaultGetDbClientDataCallBack(db)
		} else {
			arg.CallBack(db)
		}
		return
	}
	panic(fmt.Sprintf("the Database connect(%s) is not exist", arg.DbNameSpace))
}
