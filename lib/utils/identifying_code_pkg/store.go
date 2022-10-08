/**
* @Author:changjiang
* @Description:
* @File:store
* @Version: 1.0.0
* @Date 2021/2/22 12:47 上午
 */
package identifying_code_pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
)

// customizeRdsStore An object implementing Store interface
type CustomizeRdsStore struct {
	//RedisClient *redis.Client
	Context *base.Context
}

// customizeRdsStore implementing Set method of  Store interface
func (r *CustomizeRdsStore) Set(id string, value string)(err error) {
	err = r.Context.CacheClient.Set(context.Background(), id, value, time.Minute*10).Err()
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "auth.AuthLogin",
			"error":   err,
		})
	}
	return
}

// customizeRdsStore implementing Get method of  Store interface
func (r *CustomizeRdsStore) Get(id string, clear bool) (value string) {
	ctx := context.Background()
	val, err := r.Context.CacheClient.Get(ctx, id).Result()
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "auth.AuthLogin",
			"error":   err,
		})
		return
	}
	if !clear {
		return val
	}
	err = r.Context.CacheClient.Del(ctx, id).Err()
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"message": "auth.AuthLogin",
			"error":   err,
		})
		return
	}
	return val
}

func (r *CustomizeRdsStore) Verify(id, answer string, clear bool) bool {
	v := r.Get(id, clear)
	return v == answer
}

// 校验验证码类型
func (r *CustomizeRdsStore) FlagType(captchaType string) (err error) {
	var supportType = []string{
		"audio",
		"string",
		"math",
		"chinese",
		"digit",
		"",
	}
	var f bool
	for _, value := range supportType {
		if captchaType == value {
			f = true
			break
		}
	}
	if !f {
		err = fmt.Errorf("当前不支持您选择的验证码类型")
		r.Context.Error(map[string]interface{}{
			"err":             err,
			"IdentifyingCode": "IdentifyingCode.flagType",
		})
		return
	}
	return
}
