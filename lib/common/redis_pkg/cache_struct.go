package redis_pkg

import (
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
		Key           string        `json:"key"`       // key
		Expire        time.Duration `json:"expire"`    // 过期时间
		MicroApp      string        `json:"micro_app"` // 服务
		CacheDataType uint8         `json:"cache_data_type"`
		Desc          string        `json:"desc"` // 缓存使用场景描述
	}
)

func (r *CacheProperty) Default() (err error) {
	if r.CacheDataType == 0 {
		r.CacheDataType = CacheDataTypeHashString
	}
	return
}
