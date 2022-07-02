package task_timer

import (
	"context"
	"fmt"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/redis_pkg/anvil_redis"
	"github.com/robfig/cron/v3"
	"time"
)

type (
	//定时任务操作结构体
	TaskTimer struct {
		base.ControllerBase
		c     *cron.Cron
		Ctx   context.Context
		logIo *base.SystemOut
	}
	//定时任务执行过程方法
	TaskCallBack    func(context *base.Context, guid string)
	TaskTimerOption func(*TaskTimer)
)

func NewTaskTimer(options ...TaskTimerOption) (taskTimer *TaskTimer) {
	taskTimer = &TaskTimer{}
	for _, item := range options {
		item(taskTimer)
	}
	if taskTimer.logIo == nil {
		taskTimer.logIo = base.NewSystemOut()
	}
	return
}

func TaskTimerCtx(ctx context.Context) TaskTimerOption {
	return func(timer *TaskTimer) {
		timer.Ctx = ctx
	}
}

func TaskTimerCron(c *cron.Cron) TaskTimerOption {
	return func(timer *TaskTimer) {
		timer.c = c
	}
}

//配置定时任务启动
func (r *TaskTimer) ConfigTaskCall(format, taskName string, callBack TaskCallBack) (entryID cron.EntryID, err error) {

	lockKey := fmt.Sprintf("%s:%s", app_obj.App.AppName, taskName)
	r.logIo.SetInfoType(base.LogLevelInfo).
		SystemOutPrintf("【TASK】Name:%s,format:%s PK:%s",
			taskName,
			format,
			lockKey,
		)
	entryID, err = r.c.AddFunc(format, func() {
		contextBase, guid := base.CreateCrontabContext(r.ControllerBase)
		//分布式锁 防止重复执行定时任务
		_, _ = anvil_redis.NewRedisLock(
			anvil_redis.RedisLockDuration(10*time.Second),
			anvil_redis.RedisLockContext(contextBase),
			anvil_redis.RedisLockCtx(r.Ctx),
			anvil_redis.RedisLockAttemptsTime(100), //尝试获取锁的次数
			anvil_redis.RedisLockAttemptsInterval(30*time.Millisecond),
			anvil_redis.RedisLockUniqueKey(lockKey),
			anvil_redis.RedisLockLockKey(lockKey),
			anvil_redis.RedisLockOkHandler(func(ctx context.Context) (err error) {
				//执行定时任务
				callBack(contextBase, guid)
				return
			}),
		).Run()
	})
	return
}
