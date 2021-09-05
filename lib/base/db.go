// Package base
/**
* @Author:ChangJiang
* @Description:
* @File:db
* @Version: 1.0.0
* @Date 2021/1/30 4:42 下午
 */
package base

import (
	"context"
	"fmt"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"gorm.io/gorm"
)

type GetDbClientDataCallBack func(db *gorm.DB, dbName string) (dba *gorm.DB, err error)

const defaultNameSpace = "default"

// GetDbClientData 获取数据库连接 参数结构体
type GetDbClientData struct {
	Context     *Context
	DbNameSpace string                  `json:"db_name_space"`
	CallBack    GetDbClientDataCallBack // 获取数据库回调信息
}
type DbContextValue struct {
	DbName  string `json:"db_name"`
	TraceId string `json:"trace_id"`
}

func (r *GetDbClientData) DefaultGetDbClientDataCallBack(db *gorm.DB) (dba *gorm.DB, err error) {
	s, ctx := r.Context.GetTraceId()
	dba = db.WithContext(context.WithValue(ctx, app_obj.DbContextValueKey, DbContextValue{
		TraceId: s,
		DbName:  r.DbNameSpace,
	}))
	return

}

// GetDbClient 获取Redis操作实例
func GetDbClient(params ...*GetDbClientData) (db *gorm.DB, dbName string, err error) {
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
		dbName = arg.DbNameSpace
		if arg.CallBack == nil {
			db, _ = arg.DefaultGetDbClientDataCallBack(db)
		} else {
			db, _ = arg.CallBack(db, arg.DbNameSpace)
		}
		return
	} else if arg.DbNameSpace != "defaultNameSpace" { // 默认数据库连接 没有也不报错
		err = fmt.Errorf("the Database connect(%s) is not exist", arg.DbNameSpace)
		return
	}
	return
}
