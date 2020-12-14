// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package block_cache_impl

import "github.com/juetun/base-wrapper/lib/base/page_block/inte"

//TODO 缓存数据到数据库实现
func NewBlockCacheDatabaseImpl(blockCacheRedisImplOption ...BlockCacheRedisImplOption) inte.BlockCacheInterface {
	res := &blockCacheRedisImpl{}
	for _, handler := range blockCacheRedisImplOption {
		handler(res)
	}
	//初始化默认值
	res.DefaultValue()
	return res
}