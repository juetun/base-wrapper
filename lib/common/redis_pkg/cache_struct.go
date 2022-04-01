package redis_pkg

import "time"

type (
	CacheProperty struct {
		Key    string        `json:"key"`    // key
		Expire time.Duration `json:"expire"` // 过期时间
	}
)

