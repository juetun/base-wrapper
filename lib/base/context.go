// Package base
/**
* @Author:changjiang
* @Description:
* @File:Context
* @Version: 1.0.0
* @Date 2020/3/28 10:18 上午
 */
package base

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"gorm.io/gorm"
)

func CreateContext(controller *ControllerBase, c *gin.Context) (res *Context) {
	return NewContext(Log(controller.Log), GinContext(c))
}
func NewContext(contextOption ...ContextOption) *Context {
	context := &Context{}
	// 为数据指定初始化值
	for _, option := range contextOption {
		option(context)
	}
	// 初始化默认值 为空数据初始化值
	context.InitContext()
	return context
}

type Context struct {
	log         *app_obj.AppLog `json:"log"`
	Db          *gorm.DB        `json:"db"`
	CacheClient *redis.Client   `json:"cache_client"`
	GinContext  *gin.Context
	syncLog     sync.Mutex
}
type ContextOption func(context *Context)

func Log(opt *app_obj.AppLog) ContextOption {
	return func(context *Context) {
		context.log = opt
	}
}
func Db(opt *gorm.DB) ContextOption {
	return func(context *Context) {
		context.Db = opt
	}
}
func CacheClient(opt *redis.Client) ContextOption {
	return func(context *Context) {
		context.CacheClient = opt
	}
}
func GinContext(opt *gin.Context) ContextOption {
	return func(context *Context) {
		context.GinContext = opt
	}
}
func (r *Context) InitContext() (c *Context) {
	r.syncLog.Lock()
	defer r.syncLog.Unlock()
	if r.log == nil {
		r.log = app_obj.GetLog()
	}
	s := ""
	var ctx context.Context
	if nil != r.GinContext {
		if tp, ok := r.GinContext.Get(app_obj.TraceId); ok {
			s = fmt.Sprintf("%v", tp)
		}
		ctx = r.GinContext.Request.Context()
	} else {
		ctx = context.TODO()
	}
	if r.Db == nil {
		r.Db = GetDbClient(&GetDbClientData{
			Context: r,
			CallBack: func(db *gorm.DB) (dba *gorm.DB, err error) {
				dba = db.WithContext(context.WithValue(ctx, app_obj.TraceId, s))
				return
			},
		})
	}
	if r.CacheClient == nil {
		r.CacheClient = app_obj.GetRedisClient()
	}
	return r
}

func (r *Context) Error(data map[string]interface{}, message ...interface{}) {
	r.log.Error(r.GinContext, data, message)
}
func (r *Context) Info(data map[string]interface{}, message ...interface{}) {
	r.log.Info(r.GinContext, data, message)
}
func (r *Context) Debug(data map[string]interface{}, message ...interface{}) {
	r.log.Debug(r.GinContext, data, message)
}
func (r *Context) Fatal(data map[string]interface{}, message ...interface{}) {
	r.log.Fatal(r.GinContext, data, message)
}
func (r *Context) Warn(data map[string]interface{}, message ...interface{}) {
	r.log.Warn(r.GinContext, data, message)
}
