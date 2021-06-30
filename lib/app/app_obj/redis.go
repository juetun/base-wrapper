// Package app_obj
/**
* @Author:changjiang
* @Description:
* @File:db
* @Version: 1.0.0
* @Date 2020/3/27 10:39 下午
 */
package app_obj

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

var DbRedis = make(map[string]*redis.Client)

// GetRedisClient 获取Redis操作实例
func GetRedisClient(nameSpace ...string) (client *redis.Client, nameKey string) {

	switch l := len(nameSpace); l {
	case 0:
		nameKey = "default"
	case 1:
		nameKey = nameSpace[0]
	default:
		panic("nameSpace receive at most one parameter")
	}
	if _, ok := DbRedis[nameKey]; ok {
		client = DbRedis[nameKey]
		return
	}
	panic(fmt.Sprintf("the Redis connect(%s) is not exist", nameKey))
}
