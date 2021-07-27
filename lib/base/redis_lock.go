package base

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
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

	if ok, err = r.Context.CacheClient.
		SetNX(r.Context.GinContext.Request.Context(),
			r.LockKey,
			r.UniqueKey,
			r.Duration).
		Result(); err != nil {
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
	if uniqueKey != r.UniqueKey {
		err = fmt.Errorf("不是当前操作锁定数据(lock:%s,now:%s),没权限解锁", uniqueKey, r.UniqueKey)
		return
	}
	if err = r.Context.CacheClient.Del(ctx, r.LockKey).Err(); err != nil {
		r.Context.Error(map[string]interface{}{
			"err":            err.Error(),
			"LockKey":        r.LockKey,
			"redisUniqueKey": uniqueKey,
			"UniqueKey":      r.UniqueKey,
			"Duration":       r.Duration,
		}, "RedisDistributedUnLock")
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
func (r *RedisDistributedLock) tTlTime(ctx context.Context) (err error) {
	// 如果加锁成功
	// 创建协程,定时延期锁的过期时间
	for {
		select {
		case <-ctx.Done():
			// log.Printf("结束")
			return
		case <-time.After(r.Duration / 2):
			// log.Printf("续租数据\n")
			if _, err = r.addTimeout(); err != nil {
				r.Context.Error(map[string]interface{}{
					"LockKey": "续租数据",
					"err":     err.Error(),
				}, "RedisDistributedLockRun0")
			}
		}
	}
}
func (r *RedisDistributedLock) Run() (getLock bool, err error) {

	// 如果锁成功了，则操作，然后释放锁
	if getLock, err = r.Lock(); err != nil {
		return
	}

	if getLock {
		// 如果是当前操作锁定的数据
		defer func() {
			var e1 error
			if _, e1 = r.UnLock(); e1 != nil {
				r.Context.Error(map[string]interface{}{
					"err": e1.Error(),
				}, "RedisDistributedLockRun1")
			}
			// log.Println("解锁")
		}()
		ctx, cancel := context.WithCancel(r.Context.GinContext.Request.Context())
		defer func() {
			cancel()
		}()
		go func() {
			_ = r.tTlTime(ctx)
		}()

		// 执行锁逻辑
		if err = r.OkHandler(); err != nil {
			r.Context.Error(map[string]interface{}{
				"err": err.Error(),
			}, "RedisDistributedLockRun2")
			return
		}
		return
	}
	return
}

// addTimeout 延期,应该判断value后再延期
func (r *RedisDistributedLock) addTimeout() (ok bool, err error) {

	ctx := r.Context.GinContext.Request.Context()

	// 获取key的剩余有效时间 当key不存在时返回-2 当未设置过期时间的返回-1
	var ttlTime int64
	if ttlTime, err = r.Context.CacheClient.Do(ctx, "TTL", r.LockKey).Int64(); err != nil {
		r.Context.Error(map[string]interface{}{
			"LockKey": r.LockKey,
			"err":     fmt.Sprintf("redis get fail:%s", err.Error()),
		}, "RedisDistributedLockAddTimeout0")
		return
	}
	log.Printf("ttlTime:%d \n", ttlTime)
	if ttlTime <= 0 {
		ok = true
		return
	}

	if _, err = r.Context.
		CacheClient.SetEX(ctx, r.LockKey, r.UniqueKey, r.Duration).
		Result(); err != nil && err != redis.Nil {
		r.Context.Error(map[string]interface{}{
			"err":       err.Error(),
			"LockKey":   r.LockKey,
			"UniqueKey": r.UniqueKey,
			"Duration":  r.Duration,
		}, "RedisDistributedLockAddTimeout1")
		return
	}
	err = nil
	return
}
