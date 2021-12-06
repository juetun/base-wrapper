package base

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/common/response"
)

// PageCacheRedis 本功能的用途 缓存分页数据
type PageCacheRedis struct {
	CacheKeyName string                `json:"cache_key_name"` // 缓存分页数据的作用域key
	ActHandler   PageCacheRedisGetData `json:"-"`
	CacheClient  *redis.Client         `json:"-"`
	Context      *Context              `json:"-"`
	Ctx          context.Context       `json:"-"`
	Pager        *response.Pager       `json:"pager"` // 分页对象,上下文传参使用
}

// NewPageCacheRedis 初始化一个分页缓存数据对象
func NewPageCacheRedis(arg ...PageCacheRedisOption) (res *PageCacheRedis) {
	res = &PageCacheRedis{}
	for _, option := range arg {
		option(res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	return
}

type PageCacheRedisOption func(arg *PageCacheRedis)
type PageCacheRedisGetData func(pager *response.Pager, cacheKeyName string) (data interface{}, err error)

func (r *PageCacheRedis) Run(key string, data interface{}) (err error) {
	// 准备的参数验证
	if err = r.preValidate(key); err != nil {
		return
	}
	// 获取数据
	if err = r.get(key, data); err != nil && err != redis.Nil {
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"CacheKeyName": r.CacheKeyName,
			"key":          key,
		}, "PageCacheRedisRun0")
		return
	}
	if data, err = r.ActHandler(r.Pager, r.CacheKeyName); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"key":          key,
			"CacheKeyName": r.CacheKeyName,
		}, "PageCacheRedisRun1")
		return
	}
	if err = r.set(key, data); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"key":          key,
			"CacheKeyName": r.CacheKeyName,
		}, "PageCacheRedisRun2")
		return
	}
	return
}

func (r *PageCacheRedis) ClearCache() (err error) {
	if err = r.CacheClient.Del(r.Ctx, r.CacheKeyName).Err(); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"CacheKeyName": r.CacheKeyName,
		}, "PageCacheRedisClearCache")
		return
	}
	return
}

func (r *PageCacheRedis) preValidate(key string) (err error) {
	if r.Context == nil {
		err = fmt.Errorf("ctx is error")
		return
	}
	if r.CacheKeyName == "" {
		err = fmt.Errorf("cacheKeyName must be not null")
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"CacheKeyName": r.CacheKeyName,
		}, "PageCacheRedisClearCache")
		return
	}
	if key == "" {
		err = fmt.Errorf("key must be not null")
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"CacheKeyName": r.CacheKeyName,
		}, "PageCacheRedisClearCache")
		return
	}
	return
}

func (r *PageCacheRedis) set(key string, data interface{}) (err error) {
	if err = r.CacheClient.
		HSet(r.Ctx, r.CacheKeyName, key, data).
		Err(); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"CacheKeyName": r.CacheKeyName,
		}, "PageCacheRedisSet")
		return
	}
	return
}

func (r *PageCacheRedis) get(key string, data interface{}) (err error) {
	dt := r.CacheClient.HGet(r.Ctx, r.CacheKeyName, key)
	if err = dt.Err(); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"CacheKeyName": r.CacheKeyName,
			"key":          key,
		}, "PageCacheRedisGet0")
		return
	}
	if err = dt.Scan(data); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":          err.Error(),
			"CacheKeyName": r.CacheKeyName,
			"key":          key,
		}, "PageCacheRedisGet1")
		return
	}
	return
}

func PageCacheRedisCacheKeyName(CacheKeyName string) PageCacheRedisOption {
	return func(arg *PageCacheRedis) {
		arg.CacheKeyName = CacheKeyName
	}
}

func PageCacheRedisPager(pager *response.Pager) PageCacheRedisOption {
	return func(arg *PageCacheRedis) {
		arg.Pager = pager
	}
}

func PageCacheRedisActHandler(actHandler PageCacheRedisGetData) PageCacheRedisOption {
	return func(arg *PageCacheRedis) {
		arg.ActHandler = actHandler
	}
}

func PageCacheRedisCacheClient(cacheClient *redis.Client) PageCacheRedisOption {
	return func(arg *PageCacheRedis) {
		arg.CacheClient = cacheClient
	}
}

func PageCacheRedisCtx(ctx *Context) PageCacheRedisOption {
	return func(arg *PageCacheRedis) {
		arg.Context = ctx
	}
}
