package base

import (
	"fmt"
	"time"
)

// DistributedOkHandler redis 分布式锁实现结构体
type DistributedOkHandler func() (err error)

// RedisDistributedLock Redis 分布式锁
type RedisDistributedLock struct {
	AttemptsTime     int                  `json:"attempts_time"`     // 尝试获取锁的次数
	AttemptsInterval time.Duration        `json:"attempts_interval"` // 尝试获取锁时间间隔
	LockKey          string               `json:"lock_key"`
	UniqueKey        string               `json:"unique_key"`
	OkHandler        DistributedOkHandler `json:"-"`
	Context          *Context             `json:"-"`
	Duration         time.Duration        `json:"duration"`
}

func (r *RedisDistributedLock) Lock() (ok bool, err error) {
	if r.LockKey == "" || r.UniqueKey == "" {
		err = fmt.Errorf("lockKey or UniqueKey must be not null")
		r.Context.Error(map[string]interface{}{
			"err":       err.Error(),
			"LockKey":   r.LockKey,
			"UniqueKey": r.UniqueKey,
			"Duration":  r.Duration,
		}, "RedisDistributedLock0")
		return
	}
	if r.Duration == 0 {
		err = fmt.Errorf("duration is zero")
		r.Context.Error(map[string]interface{}{
			"err":       err.Error(),
			"LockKey":   r.LockKey,
			"UniqueKey": r.UniqueKey,
			"Duration":  r.Duration,
		}, "RedisDistributedLock1")
		return
	}
	ok, err = r.Context.CacheClient.
		SetNX(r.Context.GinContext.Request.Context(), r.LockKey, r.UniqueKey, r.Duration).
		Result()
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"err":       err.Error(),
			"LockKey":   r.LockKey,
			"UniqueKey": r.UniqueKey,
			"Duration":  r.Duration,
		}, "RedisDistributedLock1")
		return
	}
	return
}
func (r *RedisDistributedLock) UnLock() (ok bool, err error) {

	ctx := r.Context.GinContext.Request.Context()
	uniqueKey := r.Context.CacheClient.Get(ctx, r.LockKey).String()
	// 当前数据才能释放对应的锁
	if uniqueKey == r.UniqueKey {
		err = r.Context.CacheClient.Del(ctx, r.LockKey).Err()
		r.Context.Error(map[string]interface{}{
			"err":            err.Error(),
			"LockKey":        r.LockKey,
			"redisUniqueKey": uniqueKey,
			"UniqueKey":      r.UniqueKey,
			"Duration":       r.Duration,
		}, "RedisDistributedUnLock")
		return
	}

	return
}

// RunWithGetLock 多次尝试获取锁实现逻辑
func (r *RedisDistributedLock) RunWithGetLock() (err error) {
	var i = 0
	var getLock bool

	for {
		if i >= r.AttemptsTime {
			r.Context.Info(map[string]interface{}{
				"msg": fmt.Errorf("%d次尝试获取锁失败", r.AttemptsTime),
			}, "RedisDistributedLockRunWithGetLock")
			break
		}
		// 如果获取到锁成功，则不管执行结果如何 直接突出当前操作
		if getLock, err = r.Run(); getLock {
			return
		} else if err != nil {
			return
		}
		time.Sleep(r.AttemptsInterval)
		i++
	}
	return
}
func (r *RedisDistributedLock) Run() (getLock bool, err error) {

	// 如果锁成功了，则操作，然后释放锁
	if getLock, err = r.Lock(); err != nil {
		return
	}
	if getLock {
		// 如果是当前操作锁定的数据
		defer r.UnLock()
		err = r.OkHandler()
		return
	}
	return
}
