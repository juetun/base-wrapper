package base

import "fmt"

// 获取数据的常用类型定义
const (
	GetDataTypeFromDb    = "db"    // 从数据库获取
	GetDataTypeFromCache = "cache" // 从缓存获取
	GetDataTypeFromAll   = "all"   // 从缓存拿，如果没有则从数据库拿

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
	//GetDataWithStringIds interface {
	//	GetByIds(arg *ArgGetByStringIds) (res map[string]*models.Sku, err error)
	//	GetByIdsFromDb(id ...string)(data ,err error)
	//	GetByIdsFromCache(id ...string)(data ,err error)
	//	GetByIdsFromAll(id ...string)(data ,err error)
	//
	//}
	//GetDataWithNumberIds interface {
	//	GetByIds(arg *ArgGetByNumberIds) (res map[int64]*models.Sku, err error)
	//}
	//GetDataWithString struct {
	//
	//}
	//GetDataWithNumber struct {
	//
	//}

	GetDataTypeCommon struct {
		GetType      string `json:"get_type" form:"get_type"`
		RefreshCache uint8  `json:"refresh_cache" form:"refresh_cache"`
	}
)

func (r *GetDataTypeCommon) Default() (err error) {
	if r.GetType == "" { // 默认是从缓存拿，如果拿不到，则从数据库拿
		r.GetType = GetDataTypeFromAll
	}
	RefreshCacheValue := []uint8{0, 1}
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

//func(r *GetDataWithString)GetSKuBySkuIds(arg *ArgGetByStringIds) (res map[string]*models.Sku, err error) {
//
//	res = map[string]*models.Sku{}
//
//	if len(arg.Ids) == 0 {
//		return
//	}
//
//	arg.Default()
//
//	switch arg.GetType {
//	case GetDataTypeFromDb: // 从数据库获取数据
//		res, err = r.GetByIdsFromDb(arg.Ids...)
//	case GetDataTypeFromCache: // 从缓存获取数据
//		res, _, err = r.GetByIdsFromCache(arg.Ids...)
//	case GetDataTypeFromAll: // 优先从缓存获取，如果没有数据，则从数据库获取
//		res, err = r.GetByIdsFromAll(arg.Ids...)
//	default:
//		err = fmt.Errorf("当前不支持你选择的获取数据类型(%s)", arg.GetType)
//	}
//	return
//}
