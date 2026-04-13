package redis_pkg

import (
	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
	"time"
)

const (
	CacheDataTypeHashString  uint8 = iota + 1 //字符串
	CacheDataTypeHashHash                     //哈希
	CacheDataTypeHashList                     //列表
	CacheDataTypeHashSet                      //集合
	CacheDataTypeHashSortSet                  //有序集合
)

var (
	SliceCacheDataType = base.ModelItemOptions{
		{
			Label: "字符串",
			Value: CacheDataTypeHashString,
		},
		{
			Label: "哈希",
			Value: CacheDataTypeHashHash,
		},
		{
			Label: "列表",
			Value: CacheDataTypeHashList,
		},
		{
			Label: "集合",
			Value: CacheDataTypeHashSet,
		},
		{
			Label: "有序集合",
			Value: CacheDataTypeHashSortSet,
		},
	}
)

type (
	CacheProperty struct {
		Key      string        `json:"key"`       // key
		Expire   time.Duration `json:"expire"`    // 过期时间
		MicroApp string        `json:"micro_app"` // 服务

		CacheDataType uint8  `json:"cache_data_type"` // 缓存数据类型
		DelKey        string `json:"del_key"`         // 清除数据值
		Desc          string `json:"desc"`            // 缓存使用场景描述

		GetClientHandler GetCacheClientHandler `json:"-"` //获取缓存存储的节点方法,添加此方法可实现redis的存储数据打散到不同的redis存储
	}
	GetCacheClientHandler func(key string) (cacheClient *redis.Client, err error)
)

func (r *CacheProperty) Default() (err error) {
	if r.CacheDataType == 0 {
		r.CacheDataType = CacheDataTypeHashString
	}
	return
}
