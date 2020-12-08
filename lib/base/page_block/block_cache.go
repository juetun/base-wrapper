// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package page_block

import "time"

//缓存数据的接口，开发中需自定义实现逻辑
type BlockCacheInterface interface {
	//存储缓存数据
	//@param  name 缓存的kEY
	//@param  val  缓存的值
	//@param  cacheTime 缓存的时间
	//@return error
	Set(name string, val string, cacheTime time.Duration) (err error)

	//获取缓存数据
	//@param name 缓存的key
	//@return res 获取的数据值
	Get(name string) (res string, err error)
}

//缓存信息对象
type BlockCache struct {
	ExpireTime time.Time           `json:"expire_time"` //静态化时间周期(单位秒)，设置当前BLOCK的生命周期，如果父Block>0时以父Block的值为准。
	CacheType  string              `json:"cache_type"`  //当前界面缓存类型 如 file:文件缓存,redis:缓存，database:数据库缓存
	Cache      BlockCacheInterface `json:"cache"`       //当前界面缓存的相关信息
	CacheKey   string              `json:"cache_key"`
	CacheData  string              `json:"cache_data"` //解析后台生成的html代码，（写入缓存的数据内容）
}

func NewBlockCache(option ...BlockCacheOption) (res *BlockCache) {
	res = &BlockCache{}
	for _, handler := range option {
		handler(res)
	}
	return
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
func Cache(blockCacheInterface BlockCacheInterface) func(res *BlockCache) {
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
