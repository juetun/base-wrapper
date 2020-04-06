/**
* @Author:changjiang
* @Description:
* @File:jwt
* @Version: 1.0.0
* @Date 2020/3/28 6:40 下午
 */
package app_obj

import (
	"time"

	"github.com/go-redis/redis"
)

var jwtParam *JwtParam

func NewJwtParam() *JwtParam {
	return &JwtParam{}
}

func GetJwtParam() *JwtParam {
	return jwtParam
}

type JwtParam struct {
	DefaultIss      string
	DefaultAudience string
	DefaultJti      string
	SecretKey       string
	TokenKey        string
	TokenLife       time.Duration
	RedisCache      *redis.Client
}

func (jp *JwtParam) SetTokenKey(tk string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.TokenKey
		jp.TokenKey = tk
		return i
	}
}

func (jp *JwtParam) SetTokenLife(tl time.Duration) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.TokenLife
		jp.TokenLife = tl
		return i
	}
}

func (jp *JwtParam) SetDefaultIss(iss string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultIss
		jp.DefaultIss = iss
		return i
	}
}

func (jp *JwtParam) SetDefaultAudience(ad string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultAudience
		jp.DefaultAudience = ad
		return i
	}
}

func (jp *JwtParam) SetDefaultJti(jti string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultJti
		jp.DefaultJti = jti
		return i
	}
}

func (jp *JwtParam) SetDefaultSecretKey(sk string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.SecretKey
		jp.SecretKey = sk
		return i
	}
}

func (jp *JwtParam) SetRedisCache(rc *redis.Client) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.RedisCache
		jp.RedisCache = rc
		return i
	}
}

func (jp *JwtParam) JwtInit(options ...func(jp *JwtParam) interface{}) error {

	q := &JwtParam{
		DefaultJti:      "izghua",
		DefaultAudience: "zgh",
		DefaultIss:      "izghua",
		SecretKey:       "izghua",
		TokenLife:       time.Hour * time.Duration(72),
		TokenKey:        "login:token:",
		RedisCache:      GetRedisClient(),
	}
	for _, option := range options {
		option(q)
	}
	jwtParam = q
	return nil
}
