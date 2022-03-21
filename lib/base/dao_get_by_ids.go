package base

import "fmt"

// 获取数据的常用类型定义
const (
	GetDataTypeFromDb    = "db"    // 从数据库获取
	GetDataTypeFromCache = "cache" // 从缓存获取
	GetDataTypeFromAll   = "all"   // 从缓存拿，如果没有则从数据库拿

)

var (
	DefaultGetDataType = GetDataTypeFromDb
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
	ArgGetByStringIdsOption func(arg *ArgGetByStringIds)
	ArgGetByNumberIdsOption func(arg *ArgGetByNumberIds)
)

// NewArgGetByStringIds
func NewArgGetByStringIds(options ...ArgGetByStringIdsOption) (res *ArgGetByStringIds) {
	res = &ArgGetByStringIds{}
	for _, option := range options {
		option(res)
	}
	res.GetDataTypeCommon.Default()
	return
}

// ArgGetByStringIdsOptionGetDataTypeCommon
func ArgGetByStringIdsOptionGetDataTypeCommon(getDataTypeCommon GetDataTypeCommon) ArgGetByStringIdsOption {
	return func(arg *ArgGetByStringIds) {
		arg.GetDataTypeCommon = getDataTypeCommon
	}
}

//ArgGetByStringIdsOptionGetType
func ArgGetByStringIdsOptionGetType(getType string) ArgGetByStringIdsOption {
	return func(arg *ArgGetByStringIds) {
		arg.GetType = getType
	}
}

//ArgGetByStringIdsOptionRefreshCache
func ArgGetByStringIdsOptionRefreshCache(refreshCache uint8) ArgGetByStringIdsOption {
	return func(arg *ArgGetByStringIds) {
		arg.RefreshCache = refreshCache
	}
}

//ArgGetByStringIdsOptionIds
func ArgGetByStringIdsOptionIds(Ids ...string) ArgGetByStringIdsOption {
	return func(arg *ArgGetByStringIds) {
		arg.Ids = Ids
	}
}

//NewArgGetByNumberIds
func NewArgGetByNumberIds(options ...ArgGetByNumberIdsOption) (res *ArgGetByNumberIds) {
	res = &ArgGetByNumberIds{}
	for _, option := range options {
		option(res)
	}
	res.GetDataTypeCommon.Default()
	return
}

//ArgGetByNumberIdsOptionIds
func ArgGetByNumberIdsOptionIds(Ids ...int64) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.Ids = Ids
	}
}

//ArgGetByNumberIdsOptionGetDataTypeCommon
func ArgGetByNumberIdsOptionGetDataTypeCommon(getDataTypeCommon GetDataTypeCommon) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.GetDataTypeCommon = getDataTypeCommon
	}
}

//ArgGetByNumberIdsOptionGetType
func ArgGetByNumberIdsOptionGetType(getType string) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.GetType = getType
	}
}

//ArgGetByNumberIdsOptionRefreshCache
func ArgGetByNumberIdsOptionRefreshCache(refreshCache uint8) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.RefreshCache = refreshCache
	}
}

func (r *GetDataTypeCommon) Default() (err error) {
	if r.GetType == "" { // 默认是从缓存拿，如果拿不到，则从数据库拿
		r.GetType = DefaultGetDataType
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
