// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
// @Desc html页面缓存Redis实现方法
package block_cache_impl

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

//实现interface github.com/juetun/base-wrapper/lib/base/page_block/BlockCacheInterface
type BlockCacheRedisImpl struct {
	CacheClient *redis.Client
}

type BlockCacheRedisImplOption func(block *BlockCacheRedisImpl)

//初始化缓存对象
func NewBlockCacheRedisImpl(blockCacheRedisImplOption ...BlockCacheRedisImplOption) (res *BlockCacheRedisImpl) {
	res = &BlockCacheRedisImpl{}
	for _, handler := range blockCacheRedisImplOption {
		handler(res)
	}
	//初始化默认值
	res.defaultValue()
	return
}

//初始化默认值
func (b BlockCacheRedisImpl) defaultValue() {
	if b.CacheClient == nil {
		b.CacheClient = app_obj.GetRedisClient()
	}
}

func (b BlockCacheRedisImpl) Set(name string, val string, cacheTime time.Duration) (err error) {
	err = b.CacheClient.Set(name, val, cacheTime).Err()
	return
}

func (b BlockCacheRedisImpl) Get(name string) (res string, err error) {
	resData := b.CacheClient.Get(name)
	if err = resData.Err(); err != nil {
		return
	}
	res = resData.String()
	return
}
