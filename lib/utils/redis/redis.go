// Package redis
// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"time"
)

//这是一个redis的缓存操作工具结构体
type utilsRedis struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	redis       *redis.Client
	getDataFunc GetDataFunc
	Expiration  time.Duration
	err         error
}

type GetDataFunc func(key string) (res interface{})

func NewUtilsRedis(redisNameSpace ...string) (res *utilsRedis) {
	var nameSpace = "default"
	if len(redisNameSpace) > 0 {
		nameSpace = redisNameSpace[0]
	}
	res = &utilsRedis{}
	if _, ok := app_obj.DbRedis[nameSpace]; !ok {
		res.err = fmt.Errorf("redis connect is not exist")
		return
	}
	res.redis = app_obj.DbRedis[nameSpace]
	return
}

// Get data 必须为一个指针
func (r *utilsRedis) Get(data interface{}) (err error) {
	if r.err != nil {
		err = r.err
		return
	}
	if err = r.redis.GetSet(context.Background(), r.Key, data).Err(); err != nil {
		return
	}
	//如果缓存中没有数据
	if data == nil {
		//调用获取数据功能获取数据
		data = r.getDataFunc(r.Key)
		if err = r.Set(data); err != nil {
			return
		}
	}
	return
}

// Set data 必须为一个指针
func (r *utilsRedis) Set(data interface{}) (err error) {
	if r.err != nil {
		err = r.err
		return
	}
	if err = r.redis.Set(context.Background(), r.Key, data, r.Expiration).Err(); err != nil {
		return
	}
	return
}
func (r *utilsRedis) SetType(tp string) (res *utilsRedis) {
	r.Type = tp
	return r
}
func (r *utilsRedis) SetKey(Key string) (res *utilsRedis) {
	r.Key = Key
	return r
}
