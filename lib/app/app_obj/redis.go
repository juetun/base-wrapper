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

	"github.com/go-redis/redis"
)

var DbRedis = make(map[string]*redis.Client)

// 获取Redis操作实例
func GetRedisClient(nameSpace ...string) *redis.Client {

	var s string
	switch len := len(nameSpace); len {
	case 0:
		s = "default"
	case 1:
		s = nameSpace[0]
	default:
		panic("nameSpace receive at most one parameter")
	}
	if _, ok := DbRedis[s]; ok {
		return DbRedis[s]
	}
	panic(fmt.Sprintf("the Redis connect(%s) is not exist", s))
}
