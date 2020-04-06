/**
* @Author:changjiang
* @Description:
* @File:Context
* @Version: 1.0.0
* @Date 2020/3/28 10:18 上午
 */
package base

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/juetun/app-web/lib/app_log"
	"github.com/juetun/app-web/lib/app_obj"
)

type Context struct {
	Log         *app_log.AppLog
	Db          *gorm.DB
	CacheClient *redis.Client
}

func NewContext() *Context {
	return (&Context{}).Init()
}
func (r *Context) Init() (c *Context) {
	if r.Log == nil {
		r.Log = app_log.GetLog()
	}
	if r.Db == nil {
		r.Db = app_obj.GetDbClient()
	}
	if r.CacheClient == nil {
		r.CacheClient = app_obj.GetRedisClient()
	}
	return r
}
