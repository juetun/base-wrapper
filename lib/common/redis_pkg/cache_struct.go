package redis_pkg

import "time"

type (
	CacheProperty struct {
		Key      string        `json:"key"`       // key
		Expire   time.Duration `json:"expire"`    // 过期时间
		MicroApp string        `json:"micro_app"` // 服务
		Desc     string        `json:"desc"`      // 缓存使用场景描述
	}
)

