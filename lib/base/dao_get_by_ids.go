package parameters

// 获取数据的常用类型定义
const (
	GetDataTypeFromDb    = "db"    // 从数据库获取
	GetDataTypeFromCache = "cache" // 从缓存获取
	GetDataTypeFromAll   = "all"   // 从缓存拿，如果没有则从数据库拿

)

type GetDataTypeCommon struct {
	GetType string `json:"get_type" form:"get_type"`
}

func (r *GetDataTypeCommon) Default() {
	if r.GetType == "" { // 默认是从缓存拿，如果拿不到，则从数据库拿
		r.GetType = GetDataTypeFromAll
	}
}
