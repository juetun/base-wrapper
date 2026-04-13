package cache_act

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
	"runtime"
	"time"
)

type (
	CacheActionBase struct {
		Ctx         context.Context
		Context     *base.Context
		GetCacheKey GetCacheKey
	}
	GetCacheKey           func(id interface{}, expireTimeRand ...bool) (cacheKey string, duration time.Duration, client *redis.Client, err error)
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
	var (
		key      string
		duration time.Duration
		client   *redis.Client
	)
	if key, duration, client, err = r.GetCacheKey(id, expireTimeRand...); err != nil {
		return
	}

	if err = r.GetCacheClient(client).Set(r.Ctx, key, data, duration).Err(); err != nil {
		pc, file, line, ok := runtime.Caller(1)
		pcName := runtime.FuncForPC(pc).Name() //获取函数名
		r.Context.Info(map[string]interface{}{
			"id":       id,
			"data":     data,
			"key":      key,
			"err":      err.Error(),
			"duration": duration,
			"loc":      fmt.Sprintf("%v   %s   %d   %t   %s", pc, file, line, ok, pcName),
		}, "CacheActionSetToCache")
		return
	}
	return
}

func (r *CacheActionBase) GetCacheClient(clients ...*redis.Client) (client *redis.Client) {
	var clientTmp *redis.Client
	if len(clients) > 0 {
		clientTmp = clients[0]
	}
	if clientTmp == nil {
		clientTmp = r.Context.CacheClient
	}
	return
}

func (r *CacheActionBase) RemoveCacheByStringId(ids ...string) (err error) {
	var (
		l        = len(ids)
		keys     = make([]string, 0, l)
		key      string
		duration time.Duration
		client   *redis.Client
	)
	if l == 0 {
		return
	}
	for _, id := range ids {
		if key, _, client, err = r.GetCacheKey(id); err != nil {
			return
		}
		keys = append(keys, key)
	}

	if err = r.GetCacheClient(client).Del(r.Ctx, keys...).Err(); err != nil {
		r.Context.Info(map[string]interface{}{
			"ids":      ids,
			"keys":     keys,
			"duration": duration,
		}, "CacheActionBaseRemoveCacheByStringId")
		return
	}
	return
}

func (r *CacheActionBase) RemoveCacheByInterfaceId(ids ...interface{}) (err error) {
	var (
		l        = len(ids)
		keys     = make([]string, 0, l)
		key      string
		duration time.Duration
		client   *redis.Client
	)
	if l == 0 {
		return
	}
	for _, id := range ids {
		if key, _, client, err = r.GetCacheKey(id); err != nil {
			return
		}
		keys = append(keys, key)
	}

	if err = r.GetCacheClient(client).Del(r.Ctx, keys...).Err(); err != nil {
		r.Context.Info(map[string]interface{}{
			"ids":      ids,
			"keys":     keys,
			"duration": duration,
		}, "CacheActionBaseRemoveCacheByInterfaceId")
		return
	}
	return
}

func (r *CacheActionBase) RemoveCacheByNumberId(ids ...int64) (err error) {
	var (
		l        = len(ids)
		keys     = make([]string, 0, l)
		key      string
		duration time.Duration
		client   *redis.Client
	)
	if l == 0 {
		return
	}
	for _, id := range ids {
		if key, _, client, err = r.GetCacheKey(id); err != nil {
			return
		}
		keys = append(keys, key)
	}

	if err = r.GetCacheClient(client).Del(r.Ctx, keys...).Err(); err != nil {
		r.Context.Info(map[string]interface{}{
			"ids":      ids,
			"keys":     keys,
			"duration": duration,
		}, "CacheActionBaseRemoveCacheByNumberId")
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
