package anvil_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
)

// DistributedOkHandler redis 分布式锁实现结构体
type DistributedOkHandler func(ctx context.Context) (err error)
type RedisDistributedLockOption func(redisDistributedLock *RedisDistributedLock)

func NewRedisDistributedLock(options ...RedisDistributedLockOption) (res *RedisDistributedLock) {
	res = &RedisDistributedLock{}
	for _, option := range options {
		option(res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	return
}

func RedisDistributedLockAttemptsTime(attemptsTime int) RedisDistributedLockOption {
	return func(redisDistributedLock *RedisDistributedLock) {
		redisDistributedLock.AttemptsTime = attemptsTime
	}
}

func RedisDistributedLockAttemptsInterval(attemptsInterval time.Duration) RedisDistributedLockOption {
	return func(redisDistributedLock *RedisDistributedLock) {
		redisDistributedLock.AttemptsInterval = attemptsInterval
	}
}

func RedisDistributedLockLockKey(LockKey string) RedisDistributedLockOption {
	return func(redisDistributedLock *RedisDistributedLock) {
		redisDistributedLock.LockKey = LockKey
	}
}

func RedisDistributedLockUniqueKey(UniqueKey string) RedisDistributedLockOption {
	return func(redisDistributedLock *RedisDistributedLock) {
		redisDistributedLock.UniqueKey = UniqueKey
	}
}

func RedisDistributedLockOkHandler(OkHandler DistributedOkHandler) RedisDistributedLockOption {
	return func(redisDistributedLock *RedisDistributedLock) {
		redisDistributedLock.OkHandler = OkHandler
	}
}

func RedisDistributedLockContext(Context *base.Context) RedisDistributedLockOption {
	return func(redisDistributedLock *RedisDistributedLock) {
		redisDistributedLock.Context = Context
	}
}

func RedisDistributedLockCtx(Ctx context.Context) RedisDistributedLockOption {
	return func(redisDistributedLock *RedisDistributedLock) {
		redisDistributedLock.Ctx = Ctx
	}
}

func RedisDistributedLockDuration(Duration time.Duration) RedisDistributedLockOption {
	return func(redisDistributedLock *RedisDistributedLock) {
		redisDistributedLock.Duration = Duration
	}
}

// RedisDistributedLock Redis 分布式锁
type RedisDistributedLock struct {
	AttemptsTime     int                  `json:"attempts_time"`     // 尝试获取锁的次数
	AttemptsInterval time.Duration        `json:"attempts_interval"` // 尝试获取锁时间间隔
	LockKey          string               `json:"lock_key"`
	UniqueKey        string               `json:"unique_key"`
	OkHandler        DistributedOkHandler `json:"-"`
	Context          *base.Context        `json:"-"`
	Ctx              context.Context
	Duration         time.Duration `json:"duration"`
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
		SetNX(r.Ctx,
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

	uniqueKey := r.Context.CacheClient.Get(r.Ctx, r.LockKey).String()
	// 当前数据才能释放对应的锁
	if uniqueKey != r.UniqueKey {
		err = fmt.Errorf("不是当前操作锁定数据(lock:%s,now:%s),没权限解锁", uniqueKey, r.UniqueKey)
		return
	}
	if err = r.Context.CacheClient.Del(r.Ctx, r.LockKey).Err(); err != nil {
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

type unLockAct func()

func (r *RedisDistributedLock) tTlTime(ctx context.Context, unLockAct unLockAct) (err error) {
	// 如果加锁成功
	// 创建协程,定时延期锁的过期时间
	for {
		select {
		case <-ctx.Done():
			// log.Printf("结束")
			unLockAct()
			goto BreakLogic
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
BreakLogic:
	return
}
func (r *RedisDistributedLock) unLockAct() {
	var e1 error
	if _, e1 = r.UnLock(); e1 != nil {
		r.Context.Error(map[string]interface{}{
			"err": e1.Error(),
		}, "RedisDistributedLockRun1")
	}
	return
}
func (r *RedisDistributedLock) Run() (getLock bool, err error) {

	// 如果锁成功了，则操作，然后释放锁
	if getLock, err = r.Lock(); err != nil {
		return
	}

	if getLock {
		if r.Ctx == nil {
			r.Ctx = context.TODO()
		}
		// 如果是当前操作锁定的数据
		ctx, cancel := context.WithCancel(r.Ctx)
		defer func() {
			cancel()
		}()
		go func() {
			_ = r.tTlTime(ctx, r.unLockAct)
		}()

		// 执行锁逻辑
		if err = r.OkHandler(ctx); err != nil {
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

	logContent := map[string]interface{}{
		"LockKey": r.LockKey,
	}
	defer func() {
		if err != nil {
			r.Context.Error(logContent, "RedisDistributedLockAddTimeout")
			return
		}
		r.Context.Info(logContent, "RedisDistributedLockAddTimeout")

	}()
	// 获取key的剩余有效时间 当key不存在时返回-2 当未设置过期时间的返回-1
	var ttlTime int64
	if ttlTime, err = r.Context.CacheClient.Do(r.Ctx, "TTL", r.LockKey).Int64(); err != nil {
		logContent["desc"] = "CacheClientDo"
		return
	}
	logContent["ttlTime"] = ttlTime

	if ttlTime <= 0 {
		ok = true
		return
	}
	logContent["UniqueKey"] = r.UniqueKey
	logContent["Duration"] = r.Duration

	if _, err = r.Context.
		CacheClient.SetEX(r.Ctx, r.LockKey, r.UniqueKey, r.Duration).
		Result(); err != nil && err != redis.Nil {
		return
	}
	err = nil
	return
}
