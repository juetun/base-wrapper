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
type RedisLockOption func(RedisLock *RedisLock)

func NewRedisLock(options ...RedisLockOption) (res *RedisLock) {
	res = &RedisLock{}
	for _, option := range options {
		option(res)
	}
	if res.Ctx == nil {
		res.Ctx = context.TODO()
	}
	return
}

func RedisLockAttemptsTime(attemptsTime int) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.AttemptsTime = attemptsTime
	}
}

func RedisLockAttemptsInterval(attemptsInterval time.Duration) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.AttemptsInterval = attemptsInterval
	}
}

func RedisLockLockKey(LockKey string) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.LockKey = LockKey
	}
}

func RedisLockUniqueKey(UniqueKey string) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.UniqueKey = UniqueKey
	}
}

func RedisLockOkHandler(OkHandler DistributedOkHandler) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.OkHandler = OkHandler
	}
}

func RedisLockContext(Context *base.Context) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.Context = Context
	}
}

func RedisLockCtx(Ctx context.Context) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.Ctx = Ctx
	}
}

func RedisLockDuration(Duration time.Duration) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.Duration = Duration
	}
}

// RedisLock Redis 分布式锁
type RedisLock struct {
	AttemptsTime     int                  `json:"attempts_time"`     // 尝试获取锁的次数
	AttemptsInterval time.Duration        `json:"attempts_interval"` // 尝试获取锁时间间隔
	LockKey          string               `json:"lock_key"`
	UniqueKey        string               `json:"unique_key"`
	OkHandler        DistributedOkHandler `json:"-"`
	Context          *base.Context        `json:"-"`
	Ctx              context.Context
	Duration         time.Duration `json:"duration"`
}

func (r *RedisLock) Lock() (ok bool, err error) {
	if r.LockKey == "" || r.UniqueKey == "" {
		err = fmt.Errorf("lockKey or UniqueKey must be not null")
		r.Context.Debug(map[string]interface{}{
			"err":       err.Error(),
			"LockKey":   r.LockKey,
			"UniqueKey": r.UniqueKey,
			"Duration":  r.Duration,
		}, "RedisLock0")
		return
	}
	if r.Duration == 0 {
		err = fmt.Errorf("duration is zero")
		r.Context.Debug(map[string]interface{}{
			"err":       err.Error(),
			"LockKey":   r.LockKey,
			"UniqueKey": r.UniqueKey,
			"Duration":  r.Duration,
		}, "RedisLock1")
		return
	}

	if ok, err = r.Context.CacheClient.
		SetNX(r.Ctx,
			r.LockKey,
			r.UniqueKey,
			r.Duration).
		Result(); err != nil {
		r.Context.Debug(map[string]interface{}{
			"err":       err.Error(),
			"LockKey":   r.LockKey,
			"UniqueKey": r.UniqueKey,
			"Duration":  r.Duration,
		}, "RedisLock1")
		return
	}

	return
}
func (r *RedisLock) UnLock() (ok bool, err error) {

	uniqueKey := r.Context.CacheClient.Get(r.Ctx, r.LockKey).String()
	// 当前数据才能释放对应的锁
	if uniqueKey != r.UniqueKey {
		err = fmt.Errorf("不是当前操作锁定数据(lock:%s,now:%s),没权限解锁", uniqueKey, r.UniqueKey)
		return
	}
	if err = r.Context.CacheClient.Del(r.Ctx, r.LockKey).Err(); err != nil {
		r.Context.Debug(map[string]interface{}{
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
func (r *RedisLock) RunWithGetLock() (err error) {
	var i = 0
	var getLock bool

	for {
		if r.AttemptsTime > 0 && i >= r.AttemptsTime {
			err = fmt.Errorf("%d次尝试获取锁失败", r.AttemptsTime)
			r.Context.Warn(map[string]interface{}{
				"msg": err.Error(),
			}, "RedisLockRunWithGetLock")
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

type unLockAct func() (err error)

func (r *RedisLock) tTlTime(ctx context.Context, unLockAct unLockAct) (err error) {
	// 如果加锁成功
	// 创建协程,定时延期锁的过期时间
	for {
		select {
		case <-ctx.Done():
			// log.Printf("结束")
			_ = unLockAct()
			goto BreakLogic
		case <-time.After(r.Duration / 2):
			// log.Printf("续租数据\n")
			if _, err = r.addTimeout(); err != nil {
				r.Context.Debug(map[string]interface{}{
					"LockKey": "续租数据",
					"err":     err.Error(),
				}, "RedisLockRun0")
			}
		}
	}
BreakLogic:
	return
}
func (r *RedisLock) unLockAct() (e1 error) {
	if _, e1 = r.UnLock(); e1 != nil {
		r.Context.Debug(map[string]interface{}{
			"e1": e1.Error(),
		}, "RedisLockUnLockAct")
		return
	}
	return
}
func (r *RedisLock) Run() (getLock bool, err error) {
	logContent := map[string]interface{}{}
	defer func() {
		if err == nil {
			return
		}
		r.Context.Error(logContent, "RedisLockRun1")
	}()
	// 如果锁成功了，则操作，然后释放锁
	if getLock, err = r.Lock(); err != nil {
		logContent["desc"] = "获取锁异常"
		return
	}

	if !getLock {
		logContent["desc"] = "获取锁失败"
		return
	}
	if r.Ctx == nil {
		r.Ctx = context.TODO()
	}
	// 如果是当前操作锁定的数据
	ctx, cancel := context.WithCancel(r.Ctx)
	defer func() {
		cancel()
	}()

	go func() {
		//续租锁
		_ = r.tTlTime(ctx, r.unLockAct)
	}()

	// 执行锁逻辑
	if err = r.OkHandler(ctx); err != nil {
		r.Context.Error(map[string]interface{}{
			"err": err.Error(),
		}, "RedisLockRun2")
		return
	}

	return
}

// addTimeout 延期,应该判断value后再延期
func (r *RedisLock) addTimeout() (ok bool, err error) {

	logContent := map[string]interface{}{
		"LockKey": r.LockKey,
	}
	defer func() {
		if err != nil {
			r.Context.Error(logContent, "RedisLockAddTimeout")
			return
		}
		r.Context.Info(logContent, "RedisLockAddTimeout")

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
