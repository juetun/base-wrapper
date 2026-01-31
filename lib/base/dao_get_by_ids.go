package base

import "fmt"

// 获取数据的常用类型定义
const (
	GetDataTypeFromDb    = "db"    // 从数据库获取
	GetDataTypeFromCache = "cache" // 从缓存获取
	GetDataTypeFromAll   = "all"   // 从缓存拿，如果没有则从数据库拿

)

var (
	DefaultGetDataType = GetDataTypeFromAll
)

//是否刷新缓存
const (
	RefreshCacheNo = iota //
	RefreshCacheYes
)
const (
	ExpireTimeRandYes = true  //缓存有效期随机
	ExpireTimeRandNo  = false //缓存有效期不随机
)

type (
	ArgGetByStringIds struct {
		//parameters.GetDataTypeCommon
		GetDataTypeCommon
		Ids    []string        `json:"ids"`
		MapIds map[string]bool `json:"-"`
	}
	ArgGetByNumberIds struct {
		//parameters.GetDataTypeCommon
		GetDataTypeCommon
		Ids    []int64        `json:"ids"`
		MapIds map[int64]bool `json:"-"`
	}

	GetDataTypeCommon struct {
		InitCount      int    `json:"-" form:"-"`
		GetType        string `json:"get_type" form:"get_type"`                 // 本次读取数据从数据库 ，缓存，或使用优先从缓存(缓存没有,则从数据库读取 同时写入缓存)
		IncludeDelData bool   `json:"include_del_data" form:"include_del_data"` // 查询数据包括已删除(软删)的数据默认查询不包括软删数据
		RefreshCache   uint8  `json:"refresh_cache" form:"refresh_cache"`       // 是否刷新缓存数据
		MaxLimit       int64  `json:"max_limit" form:"max_limit"`               // 本次请求最多查询数据数量
		ExpireTimeRand bool   `json:"expire_time_rand" form:"expire_time_rand"` // 缓存有效期是否
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

	_ = res.GetDataTypeCommon.Default()
	res.Ids = make([]string, 0, res.InitCount)
	res.MapIds = make(map[string]bool, res.InitCount)
	return
}

func (r *ArgGetByNumberIds) AddIds(ids ...int64) {
	var (
		ok bool
		id int64
	)
	for _, id = range ids {
		if _, ok = r.MapIds[id]; !ok {
			r.Ids = append(r.Ids, id)
			r.MapIds[id] = true
		}
	}
	return
}

func (r *ArgGetByNumberIds) GetIdMap() (res map[int64]bool) {
	res = make(map[int64]bool, len(r.Ids))
	for _, item := range r.Ids {
		res[item] = true
	}
	return
}

func (r *ArgGetByStringIds) GetIdMap() (res map[string]bool) {
	res = make(map[string]bool, len(r.Ids))
	for _, item := range r.Ids {
		res[item] = true
	}
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

func ArgGetByStringIdsOptionExpireTimeRand(expireTimeRand bool) ArgGetByStringIdsOption {
	return func(arg *ArgGetByStringIds) {
		arg.ExpireTimeRand = expireTimeRand
	}
}

func ArgGetByStringIdsOptionInitCount(initCount int) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.InitCount = initCount
	}
}

func (r *ArgGetByStringIds) AddIds(ids ...string) {
	var (
		ok bool
		id string
	)
	for _, id = range ids {
		if _, ok = r.MapIds[id]; !ok {
			r.Ids = append(r.Ids, id)
			r.MapIds[id] = true
		}
	}
	return
}

//NewArgGetByNumberIds
func NewArgGetByNumberIds(options ...ArgGetByNumberIdsOption) (res *ArgGetByNumberIds) {
	res = &ArgGetByNumberIds{}
	for _, option := range options {
		option(res)
	}
	_ = res.GetDataTypeCommon.Default()
	res.Ids = make([]int64, 0, res.InitCount)
	res.MapIds = make(map[int64]bool, res.InitCount)
	return
}

//ArgGetByNumberIdsOptionIds
func ArgGetByNumberIdsOptionIds(Ids ...int64) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.Ids = Ids
	}
}
func ArgGetByNumberIdsOptionExpireTimeRand(expireTimeRand bool) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.ExpireTimeRand = expireTimeRand
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

func ArgGetByNumberIdsOptionInitCount(initCount int) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.InitCount = initCount
	}
}

//ArgGetByNumberIdsOptionRefreshCache
func ArgGetByNumberIdsOptionRefreshCache(refreshCache uint8) ArgGetByNumberIdsOption {
	return func(arg *ArgGetByNumberIds) {
		arg.RefreshCache = refreshCache
	}
}
func (r *GetDataTypeCommon) validateExpireTimeRandValue() (err error) {
	ExpireTimeRandValue := []bool{ExpireTimeRandYes, ExpireTimeRandNo}
	var f1 bool
	for _, value := range ExpireTimeRandValue {
		if value == r.ExpireTimeRand {
			f1 = true
			break
		}
	}
	if !f1 {
		err = fmt.Errorf("expire_time_rand值异常")
		return
	}
	return
}

func (r *GetDataTypeCommon) validateRefreshCacheValue() (err error) {
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

func (r *GetDataTypeCommon) Default() (err error) {
	switch r.RefreshCache {
	case RefreshCacheYes:
		r.GetType = GetDataTypeFromDb
	}
	if r.GetType == "" { //如果没有设置获取数据的类型
		r.GetType = DefaultGetDataType // 默认是从缓存拿，如果拿不到，则从数据库拿
	}
	if err = r.validateExpireTimeRandValue(); err != nil {
		return
	}
	if err = r.validateRefreshCacheValue(); err != nil {
		return
	}

	return
}
