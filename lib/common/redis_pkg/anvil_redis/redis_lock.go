package anvil_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/base"
)

//分布式锁，锁数据
/** 调用实例
_ = anvil_redis.NewRedisLock(
		anvil_redis.RedisLockDuration(10*time.Second),
		anvil_redis.RedisLockContext(r.Context),
		anvil_redis.RedisLockCtx(r.Ctx),
		anvil_redis.RedisLockAttemptsTime(100), //尝试获取锁的次数
		anvil_redis.RedisLockAttemptsInterval(30*time.Millisecond),
		anvil_redis.RedisLockUniqueKey(utils.Guid(uk)),
		anvil_redis.RedisLockLockKey(uk),
		anvil_redis.RedisLockOkHandler(func(ctx context.Context) (err error) {
			// TODO 获取到锁的逻辑 ...
			return
		}),
	).RunWithGetLock()
*/
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

// DistributedOkHandler redis 分布式锁实现结构体
type (
	RedisLock struct {
		AttemptsTime     int                  `json:"attempts_time"`     // 尝试获取锁的次数
		AttemptsInterval time.Duration        `json:"attempts_interval"` // 尝试获取锁时间间隔
		LockKey          string               `json:"lock_key"`          // 业务key
		UniqueKey        string               `json:"unique_key"`        // 锁的值 （用于释放锁时，只能指定线程才可释放锁）
		expiration       time.Duration        `json:"duration"`          // 锁的时长 （单位秒）
		tTlDuration      time.Duration        `json:"ttl_duration"`      // 锁续时的时长 （单位秒）
		maxExecDuration  time.Duration        `json:"max_exec_duration"` //最大执行时长 0表示无限制
		RedisClient      *redis.Client        `json:"-"`
		OkHandler        DistributedOkHandler `json:"-"`
		Context          *base.Context        `json:"-"`
		Ctx              context.Context      `json:"-"`
	}

	DistributedOkHandler func(ctx context.Context) (err error)
	RedisLockOption      func(RedisLock *RedisLock)
	unLockAct            func() (err error)
)

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
func RedisLockRedisClient(redisClient *redis.Client) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.RedisClient = redisClient
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

//key的生命周期
func RedisLockDuration(Duration time.Duration) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.expiration = Duration
	}
}

//
func RedisLockMaxExecDuration(maxExecDuration time.Duration) RedisLockOption {
	return func(RedisLock *RedisLock) {
		RedisLock.maxExecDuration = maxExecDuration
	}
}

// RedisLock Redis 分布式锁
// RunWithGetLock 多次尝试获取锁实现逻辑
func (r *RedisLock) RunWithGetLock() (err error) {
	//校验参数是否可用
	if err = r.validateParameters(); err != nil {
		return
	}
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

//单次去获取锁，获取到就做 没获取到就跳过
//	options := []anvil_redis.RedisLockOption{anvil_redis.RedisLockDuration(10 * time.Second),
//	anvil_redis.RedisLockContext(r.Context),
//	anvil_redis.RedisLockCtx(r.ctx),
//	anvil_redis.RedisLockRedisClient(client),
//	anvil_redis.RedisLockDuration(3 * time.Second),
//	anvil_redis.RedisLockUniqueKey(utils.Guid(cacheKey)),
//	anvil_redis.RedisLockLockKey(cacheKey),
//	anvil_redis.RedisLockOkHandler(func(ctx context.Context) (err error) {
//	_ = r.logicLoadData(client, cacheKey, duration, status, timeNow, notLoadingCallBack)
//	return
//	}),}
//	if _, err = anvil_redis.NewRedisLock(options...).
//	Run(); err != nil {
//	return
//	}
func (r *RedisLock) Run() (getLock bool, err error) {
	if err = r.validateParameters(); err != nil {
		return
	}
	logContent := map[string]interface{}{
		"lock_key":   r.LockKey,
		"unique_key": r.UniqueKey,
	}
	defer func() {
		mark := "RedisLockRun1"
		if err == nil {
			r.Context.Debug(logContent, mark)
			return
		}
		logContent["err"] = err.Error()
		r.Context.Error(logContent, mark)
	}()
	// 如果锁成功了，则操作，然后释放锁
	if getLock, err = r.lock(); err != nil {
		logContent["desc"] = "获取锁异常"
		return
	}

	if !getLock {
		logContent["desc"] = "获取锁失败"
		return
	}
	var cancel context.CancelFunc
	if r.maxExecDuration > 0 { //如果设置了最大执行时长
		// 如果是当前操作锁定的数据
		r.Ctx, cancel = context.WithTimeout(r.getCtx(), r.maxExecDuration)
	} else {
		// 如果是当前操作锁定的数据
		r.Ctx, cancel = context.WithCancel(r.getCtx())
	}

	go func() {
		defer cancel()
		// 执行锁逻辑
		if err = r.OkHandler(r.Ctx); err != nil {
			r.Context.Error(map[string]interface{}{
				"err": err.Error(),
			}, "RedisLockRun2")
			return
		}

	}()

	//续租锁
	r.ttlRun(r.Ctx, r.unLockAct)
	return
}

func (r *RedisLock) lock() (ok bool, err error) {
	defer func() {
		if err == nil {
			return
		}
		r.Context.Debug(map[string]interface{}{
			"err":        err.Error(),
			"LockKey":    r.LockKey,
			"UniqueKey":  r.UniqueKey,
			"expiration": r.expiration,
		}, "RedisLock")
	}()
	if r.LockKey == "" || r.UniqueKey == "" {
		err = fmt.Errorf("lockKey or UniqueKey must be not null")
		return
	}
	if r.expiration == 0 {
		err = fmt.Errorf("duration is zero")
		return
	}

	if ok, err = r.getCacheClient().
		SetNX(r.getCtx(),
			r.LockKey,
			r.UniqueKey,
			r.expiration).
		Result(); err != nil {
		return
	}
	return
}

func (r *RedisLock) unLock() (ok bool, err error) {

	script := redis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`)

	//
	ctx := context.TODO() //注(请勿修改):此上下文句柄为操作Redis使用 不要使用 锁数据的句柄
	result, err := script.Run(ctx, r.getCacheClient(), []string{r.LockKey}, r.UniqueKey).Int64()
	if err != nil {
		return
	}
	ok = result > 0
	return
}

func (r *RedisLock) validateParameters() (err error) {
	if r.expiration == 0 { //锁的有效期
		err = fmt.Errorf("请设置Duration")
		return
	}
	if r.tTlDuration >= r.expiration {
		err = fmt.Errorf("Duration必须大于TTlDuration")
		return
	}
	if r.tTlDuration == 0 {
		r.tTlDuration = r.expiration - 1
		if r.tTlDuration <= 0 {
			err = fmt.Errorf("请设置TTlDuration的值")
			return
		}
	}
	return
}

func (r *RedisLock) ttlRun(ctx context.Context, unLockAct unLockAct) () {
	var err error
	// 如果加锁成功
	// 创建协程,定时延期锁的过期时间
	for {
		select {
		case <-ctx.Done():
			// log.Printf("结束")
			_ = unLockAct()
			goto BreakLogic
		case <-time.After(r.tTlDuration):
			mark := "RedisLockRun0"
			if _, err = r.refreshLock(); err != nil {
				r.Context.Error(map[string]interface{}{
					"LockKey": "续租失败",
					"err":     err.Error(),
				}, mark)
			} else {
				r.Context.Debug(map[string]interface{}{
					"LockKey": "续租数据",
				}, mark)
			}
		}
	}
BreakLogic:
	return
}

func (r *RedisLock) unLockAct() (e1 error) {
	if _, e1 = r.unLock(); e1 != nil {
		r.Context.Debug(map[string]interface{}{
			"e1": e1.Error(),
		}, "RedisLockUnLockAct")
		return
	}
	return
}

func (r *RedisLock) getCacheClient() (res *redis.Client) {
	if r.RedisClient != nil {
		res = r.RedisClient
		return
	}
	res = r.Context.CacheClient
	return
}

func (r *RedisLock) getCtx() (res context.Context) {
	if r.Ctx != nil {
		res = r.Ctx
		return
	}
	res = context.TODO()
	return
}

// RefreshLock 存在则更新过期时间,不存在则创建key
func (r *RedisLock) refreshLock() (ok bool, err error) {
	script := redis.NewScript(`
	local val = redis.call("GET", KEYS[1])
	if not val then
		redis.call("setex", KEYS[1], ARGV[2], ARGV[1])
		return 2
	elseif val == ARGV[1] then
		return redis.call("expire", KEYS[1], ARGV[2])
	else
		return 0
	end
	`)
	var result int64
	if result, err = script.Run(r.getCtx(), r.getCacheClient(), []string{r.LockKey}, r.UniqueKey, r.expiration/time.Second).Int64(); err != nil {
		return
	} else {
		ok = result > 0
	}
	return
}
