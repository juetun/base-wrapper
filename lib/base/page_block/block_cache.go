// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package page_block

import (
	"log"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base/page_block/block_cache_impl"
	"github.com/juetun/base-wrapper/lib/base/page_block/inte"
)

const (
	CacheFile     = "file"     //缓存到文件
	CacheRedis    = "redis"    //缓存到Redis
	CacheDatabase = "database" //缓存到数据库
)

//缓存信息对象
type BlockCache struct {
	ExpireTime time.Time                `json:"expire_time"` //静态化时间周期(单位秒)，设置当前BLOCK的生命周期，如果父Block>0时以父Block的值为准。
	CacheType  string                   `json:"cache_type"`  //当前界面缓存类型 如 file:文件缓存,redis:缓存，database:数据库缓存
	Cache      inte.BlockCacheInterface `json:"cache"`       //当前界面缓存的相关信息
	CacheKey   string                   `json:"cache_key"`
	CacheData  string                   `json:"cache_data"` //解析后台生成的html代码，（写入缓存的数据内容）
}

func NewBlockCache(option ...BlockCacheOption) (res *BlockCache) {
	res = &BlockCache{
		CacheType: CacheRedis,
	}
	for _, handler := range option {
		handler(res)
	}

	//初始化默认值
	res.Default()
	return
}
func (r *BlockCache) Default() {
	if r.Cache == nil {
		r.defaultCache()
	}

}
func (r *BlockCache) defaultCache() {
	switch strings.ToLower(r.CacheType) {
	case CacheRedis: //缓存到redis
		r.Cache = block_cache_impl.NewBlockCacheRedisImpl()
	case CacheFile: //缓存到文件
		r.Cache = block_cache_impl.NewBlockCacheFileImpl()
	case CacheDatabase: //缓存到数据库
		r.Cache = block_cache_impl.NewBlockCacheDatabaseImpl()
	default:
		log.Fatalf("the cache type is not supported (%s)", r.CacheType)
	}
}

type BlockCacheOption func(block *BlockCache)

func CacheData(cacheData string) func(res *BlockCache) {
	return func(blockCache *BlockCache) {
		blockCache.CacheData = cacheData
	}
}
func CacheKey(cacheKey string) func(res *BlockCache) {
	return func(blockCache *BlockCache) {
		blockCache.CacheKey = cacheKey
	}
}
func Cache(blockCacheInterface inte.BlockCacheInterface) func(res *BlockCache) {
	return func(blockCache *BlockCache) {
		blockCache.Cache = blockCacheInterface
	}
}
func CacheType(cacheType string) func(blockCache *BlockCache) {
	return func(blockCache *BlockCache) {
		blockCache.CacheType = cacheType
	}
}

func ExpireTime(tt time.Time) func(blockCache *BlockCache) {
	return func(blockCache *BlockCache) {
		blockCache.ExpireTime = tt
	}
}
