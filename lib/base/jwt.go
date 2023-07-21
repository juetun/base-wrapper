// Package app_obj
/**
* @Author:changjiang
* @Description:
* @File:jwt
* @Version: 1.0.0
* @Date 2020/3/28 6:40 下午
 */
package base

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

// 当前请求上下文存储使用的KEY
const (
	ContextUserObjectKey = "jwt_user" // 用户信息
	ContextUserTokenKey  = "token"    // 存储token的KEY
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
type JwtUser struct {
	UserId int64  `json:"user_hid"` // 用户ID
	Name   string `json:"name"`     // 用户昵称
	// Portrait string `json:"portrait"` // 头像
	Status int8 `json:"status"` // '用户状态 0创建,1正常',
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
	cacheClient, _ := app_obj.GetRedisClient()
	q := &JwtParam{
		DefaultJti:      "juetun",
		DefaultAudience: "zgh",
		DefaultIss:      "juetun",
		SecretKey:       "juetun",
		TokenLife:       time.Hour * time.Duration(72),
		TokenKey:        "log_tk:",
		RedisCache:      cacheClient,
	}
	for _, option := range options {
		option(q)
	}
	jwtParam = q
	return nil
}

func CreateTokenFromObject(user interface{}, c *Context) (tokenString string, err error) {
	var s []byte
	if s, err = json.Marshal(user); err != nil {
		return
	}
	if tokenString, err = CreateToken(string(s), c); err != nil {
		return
	}
	return
}

func CreateToken(s string, ctx *Context) (tokenString string, err error) {
	logContent := make(map[string]interface{}, 10)
	defer func() {
		if err == nil {
			return
		}
		ctx.Error(logContent, "jwtCreateToken")
	}()

	//	iss: jwt签发者
	//	sub: jwt所面向的用户
	//	aud: 接收jwt的一方
	//	exp: jwt的过期时间，这个过期时间必须要大于签发时间
	//	nbf: 定义在什么时间之前，该jwt都是不可用的.
	//	iat: jwt的签发时间
	//	jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
	jwtParam := GetJwtParam()
	tk := jwt.New(jwt.SigningMethodHS256)
	var claims jwt.MapClaims = map[string]interface{}{
		"iat": time.Now().Unix(),
		"iss": jwtParam.DefaultIss,
		"sub": s,
		"aud": jwtParam.DefaultAudience,
		"jti": utils.Md5(jwtParam.DefaultJti + jwtParam.DefaultIss),
	}

	tk.Claims = claims
	SecretKey := jwtParam.SecretKey
	if tokenString, err = tk.SignedString([]byte(SecretKey)); err != nil {

		logContent["content"] = "token create error"
		logContent["error"] = err.Error()

		return
	}
	if jwtParam.RedisCache == nil {
		err = fmt.Errorf("redis connect is not exists")
		logContent["content"] = "common/jwt.go"
		logContent["error"] = err.Error()
		return
	}
	err = jwtParam.RedisCache.
		Set(ctx.GinContext.Request.Context(), jwtParam.TokenKey+tokenString, s, jwtParam.TokenLife).
		Err()
	if err != nil {
		err = fmt.Errorf("redis connect is not exists")
		logContent["content"] = "token create error"
		logContent["error"] = err.Error()
		return
	}

	return
}

func ParseToken(myToken string, ctx *Context) (sub string, err error) {
	if myToken == "" || myToken == "null" {
		return
	}
	var (
		logContent = make(map[string]interface{}, 6)
		jwtParam   = GetJwtParam()
		token      *jwt.Token
		ok         bool
		res        string
	)
	defer func() {
		if err == nil || ctx == nil {
			return
		}
		ctx.Error(logContent, "jwtParseToken")
	}()
	logContent["token"] = "'" + myToken + "'"
	defer func() {
		if err != nil {
			logContent["error"] = err.Error()
		}
	}()
	if myToken == "" {
		logContent["content"] = "token is null"
		return
	}

	if token, err = jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	}); err != nil {
		logContent["content"] = "parse token has error"
		return
	}

	if !token.Valid {
		logContent["content"] = fmt.Sprintf("token is invalid(%s)", myToken)
		return
	}
	var (
		claims = token.Claims.(jwt.MapClaims)
	)
	if sub, ok = claims["sub"].(string); !ok {
		logContent["content"] = "claimsSub"
		err = fmt.Errorf("claims duan yan is error")
		return
	}
	if jwtParam.RedisCache == nil {
		logContent["content"] = "common/jwt.go"
		err = fmt.Errorf("Redis connect is null ")
		return
	}
	cacheKey := jwtParam.TokenKey + myToken

	var e error
	//检测缓存中是否有token
	res, e = jwtParam.RedisCache.
		Get(ctx.GinContext.Request.Context(), cacheKey).
		Result()
	if e == redis.Nil { //如果缓存中没有数据
		return
	}
	if err = e; err != nil {
		logContent["content"] = "get token from redis error"
		return
	}

	if res == "" {
		desc := "token is invalid"
		logContent["content"] = desc
		err = fmt.Errorf(desc)
		return
	}

	// refresh the token life time
	if err = jwtParam.RedisCache.Set(ctx.GinContext.Request.Context(), cacheKey, myToken, jwtParam.TokenLife).Err(); err != nil {
		logContent["content"] = "token create error"
		return
	}
	return
}

//解析用户JWT key
func ParseJwtKey(myToken string, ctx *Context, data interface{}) (err error) {
	var sub string
	if sub, err = ParseToken(myToken, ctx); err != nil {
		return
	}
	logContent := make(map[string]interface{}, 1)
	defer func() {
		if err == nil {
			return
		}
		ctx.Error(logContent, "jwtParseUser")
	}()
	if err = json.Unmarshal([]byte(sub), data); err != nil {
		logContent["content"] = "sub is error may be is not a json string"
		return
	}
	return
}

func UnsetToken(myToken string, ctx *Context) (ok bool, err error) {
	jwtParam := GetJwtParam()
	logContent := make(map[string]interface{}, 10)
	var (
		sub string
	)
	defer func() {
		if err == nil {
			return
		}
		ctx.Error(map[string]interface{}{
			"content": "parse token has error",
			"error":   err.Error(),
			"token":   "'" + myToken + "'",
		}, "jwtUnsetToken")
	}()
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (res interface{}, e error) {
		res = []byte(jwtParam.SecretKey)
		return
	})
	if err != nil {
		logContent ["content"] = fmt.Sprintf("parse token has error ,token:`%s`", myToken)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	if sub, ok = claims["sub"].(string); !ok {
		logContent ["content"] = "claims duan yan is error"
		err = fmt.Errorf("claims duan yan is error")
		return
	}

	if err = jwtParam.RedisCache.Del(ctx.GinContext.Request.Context(), jwtParam.TokenKey+sub).Err(); err != nil {
		logContent ["content"] = "unset token has error"
		return
	}
	ok = true
	return
}
