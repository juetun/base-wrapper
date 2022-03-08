package base

import "fmt"

// 获取数据的常用类型定义
const (
	GetDataTypeFromDb    = "db"    // 从数据库获取
	GetDataTypeFromCache = "cache" // 从缓存获取
	GetDataTypeFromAll   = "all"   // 从缓存拿，如果没有则从数据库拿

)

//是否刷新缓存
const (
	RefreshCacheNo = iota
	RefreshCacheYes
)

type (
	ArgGetByStringIds struct {
		//parameters.GetDataTypeCommon
		GetDataTypeCommon
		Ids []string `json:"ids"`
	}
	ArgGetByNumberIds struct {
		//parameters.GetDataTypeCommon
		GetDataTypeCommon
		Ids []int64 `json:"ids"`
	}
	
	GetDataTypeCommon struct {
		GetType      string `json:"get_type" form:"get_type"`
		RefreshCache uint8  `json:"refresh_cache" form:"refresh_cache"`
	}
)

func (r *GetDataTypeCommon) Default() (err error) {
	if r.GetType == "" { // 默认是从缓存拿，如果拿不到，则从数据库拿
		r.GetType = GetDataTypeFromAll
	}

	RefreshCacheValue := []uint8{RefreshCacheNo, RefreshCacheYes}
	var f bool
	for _, value := range RefreshCacheValue {
		if value == r.RefreshCache {
			f = true
			break
		}
	}
	if !f {
		err = fmt.Errorf("refresh_cache值异常")
		return
	}
	return
}
