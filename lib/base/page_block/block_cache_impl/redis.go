// Package block_cache_impl
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
// @Desc html页面缓存Redis实现方法
package block_cache_impl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base/page_block/inte"
)

// 实现interface github.com/juetun/base-wrapper/lib/base/page_block/BlockCacheInterface
type blockCacheRedisImpl struct {
	CacheClient    *redis.Client
	CacheNameSapce string
}

type BlockCacheRedisImplOption func(block inte.BlockCacheInterface)

// NewBlockCacheRedisImpl 初始化缓存对象
func NewBlockCacheRedisImpl(blockCacheRedisImplOption ...BlockCacheRedisImplOption) inte.BlockCacheInterface {
	res := &blockCacheRedisImpl{}
	for _, handler := range blockCacheRedisImplOption {
		handler(res)
	}
	// 初始化默认值
	res.DefaultValue()
	return res
}

// DefaultValue 初始化默认值
func (b *blockCacheRedisImpl) DefaultValue() {
	if b.CacheClient == nil {

		b.CacheClient, b.CacheNameSapce = app_obj.GetRedisClient()
		if b.CacheClient == nil {
			panic(fmt.Errorf("get cache client exception"))
		}

	}
	return
}

// 写入缓存数据
func (b *blockCacheRedisImpl) Set(name string, val string, cacheTime time.Duration) (err error) {

	err = b.CacheClient.Set(context.Background(), name, val, cacheTime).Err()
	return
}

// 读取缓存数据
func (b *blockCacheRedisImpl) Get(name string) (res string, err error) {
	resData := b.CacheClient.Get(context.Background(), name)
	if err = resData.Err(); err != nil {
		return
	}
	res = resData.String()
	return
}
