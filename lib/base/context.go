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
	"net/http"
	"sync"

	redis "github.com/go-redis/redis/v8"

	"github.com/juetun/base-wrapper/lib/utils"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"gorm.io/gorm"
)

func CreateCrontabContext(controllerBase ControllerBase, uniqueKey ...string) (context *Context, guid string) {
	guid = utils.Guid("crontabContext")
	if len(uniqueKey) > 0 {
		guid = uniqueKey[0]
	}
	ginContext := &gin.Context{Request: &http.Request{Header: map[string][]string{}}}
	ginContext.Request.Header.Set(app_obj.HttpTraceId, guid)
	ginContext.Set(app_obj.TraceId, guid)
	context = CreateContext(&controllerBase, ginContext)
	return
}
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
	log            *app_obj.AppLog           `json:"log"`
	Db             *gorm.DB                  `json:"db"`
	DbName         string                    `json:"db_name"` // 数据库的链接配置KEY
	CacheClient    *redis.Client             `json:"cache_client"`
	CacheName      string                    `json:"cache_name"` // 缓存的链接配置KEY
	ClickHouse     *app_obj.ClickHouseClient `json:"click_house"`
	ClickHouseName string                    `json:"click_house_name"`
	GinContext     *gin.Context
	// syncLog        sync.Mutex
}
type ContextOption func(context *Context)

func Log(opt *app_obj.AppLog) ContextOption {
	return func(context *Context) {
		context.log = opt
	}
}

// GetTraceId 根据请求获取trace_id
func (r *Context) GetTraceId(ctxObj ...context.Context) (res string, ctx context.Context) {
	res = ""
	if len(ctxObj) > 0 {
		ctx = ctxObj[0]
	} else {
		ctx = context.TODO()
	}
	if nil != r.GinContext {
		if tp, ok := r.GinContext.Get(app_obj.TraceId); ok {
			res = fmt.Sprintf("%v", tp)
		}
	}
	return
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
func (r *Context) InitContext(needLock ...bool) (c *Context) {
	if len(needLock) > 0 && needLock[0] {
		var syncLog sync.Mutex
		syncLog.Lock()
		defer syncLog.Unlock()
	}

	if r.log == nil {
		r.log = app_obj.GetLog()
	}
	s, ctx := r.GetTraceId()
	if r.Db == nil {
		_ = r.initDb(s, ctx)
	}
	if r.ClickHouse == nil {
		r.initClickHouse()
	}
	if r.CacheClient == nil {
		r.CacheClient, r.CacheName = app_obj.GetRedisClient()
	}
	return r
}

// 初始化clickhouse句柄
func (r *Context) initClickHouse() {
	r.ClickHouse, r.ClickHouseName = app_obj.GetClickHouseClient()
}

func (r *Context) initDb(s string, ctx context.Context) (err error) {
	r.Db, r.DbName, err = GetDbClient(&GetDbClientData{
		Context: r,
		CallBack: func(db *gorm.DB, dbName string) (dba *gorm.DB, err error) {
			dba = db.WithContext(context.WithValue(ctx, app_obj.DbContextValueKey, DbContextValue{
				TraceId: s,
				DbName:  dbName,
			}))
			return
		},
	})

	return
}

func (r *Context) Error(data map[string]interface{}, message ...interface{}) {
	r.log.Error(r.GinContext, data, message...)
}
func (r *Context) Info(data map[string]interface{}, message ...interface{}) {
	r.log.Info(r.GinContext, data, message...)
}
func (r *Context) Debug(data map[string]interface{}, message ...interface{}) {
	r.log.Debug(r.GinContext, data, message...)
}
func (r *Context) Fatal(data map[string]interface{}, message ...interface{}) {
	r.log.Fatal(r.GinContext, data, message...)
}
func (r *Context) Warn(data map[string]interface{}, message ...interface{}) {
	r.log.Warn(r.GinContext, data, message...)
}
