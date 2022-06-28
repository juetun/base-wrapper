package cache_act

import (
	"context"
	"github.com/juetun/base-wrapper/lib/base"
	"time"
)

type (
	CacheActionBase struct {
		Ctx         context.Context
		Context     *base.Context
		GetCacheKey GetCacheKey
	}
	GetCacheKey           func(id interface{}, expireTimeRand ...bool) (cacheKey string, duration time.Duration)
	CacheActionBaseOption func(cacheFreightAction *CacheActionBase)
)

func NewCacheActionBasePointer(options ...CacheActionBaseOption) (res *CacheActionBase) {
	res = &CacheActionBase{}
	for _, handler := range options {
		handler(res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	return
}

func NewCacheActionBase(options ...CacheActionBaseOption) (res CacheActionBase) {
	res = CacheActionBase{}
	for _, handler := range options {
		handler(&res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	return
}

func (r *CacheActionBase) LoadBase(options ...CacheActionBaseOption) {
	for _, handler := range options {
		handler(r)
	}
}

func (r *CacheActionBase) SetToCache(id interface{}, data interface{}, expireTimeRand ...bool) (err error) {
	key, duration := r.GetCacheKey(id, expireTimeRand...)
	if err = r.Context.CacheClient.Set(r.Ctx, key, data, duration).Err(); err != nil {
		r.Context.Info(map[string]interface{}{
			"id":       id,
			"data":     data,
			"key":      key,
			"duration": duration,
		}, "CacheActionSetToCache")
		return
	}
	return
}

func CacheActionBaseGetCacheKey(getCacheKey GetCacheKey) CacheActionBaseOption {
	return func(cacheFreightAction *CacheActionBase) {
		cacheFreightAction.GetCacheKey = getCacheKey
		return
	}
}

func CacheActionBaseContext(context *base.Context) CacheActionBaseOption {
	return func(cacheFreightAction *CacheActionBase) {
		cacheFreightAction.Context = context
		return
	}
}

func CacheActionBaseCtx(ctx context.Context) CacheActionBaseOption {
	return func(cacheFreightAction *CacheActionBase) {
		cacheFreightAction.Ctx = ctx
		return
	}
}
