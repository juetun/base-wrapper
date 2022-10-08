// Package common
// /**
package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	utils2 "github.com/juetun/base-wrapper/lib/utils"
)

func CreateToken(user app_obj.JwtUser, c *gin.Context) (tokenString string, err error) {
	//	iss: jwt签发者
	//	sub: jwt所面向的用户
	//	aud: 接收jwt的一方
	//	exp: jwt的过期时间，这个过期时间必须要大于签发时间
	//	nbf: 定义在什么时间之前，该jwt都是不可用的.
	//	iat: jwt的签发时间
	//	jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
	jwtParam := app_obj.GetJwtParam()
	tk := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(time.Hour * time.Duration(72)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = jwtParam.DefaultIss
	s, _ := json.Marshal(user)
	userIdString := string(s)
	claims["sub"] = userIdString
	claims["aud"] = jwtParam.DefaultAudience
	claims["jti"] = utils2.Md5(jwtParam.DefaultJti + jwtParam.DefaultIss)
	tk.Claims = claims
	SecretKey := jwtParam.SecretKey
	tokenString, err = tk.SignedString([]byte(SecretKey))
	if err != nil {
		app_obj.GetLog().Error(c, map[string]interface{}{
			"content": "token create error",
			"error":   err.Error(),
		})
		return
	}
	if jwtParam.RedisCache == nil {
		app_obj.GetLog().Error(c, map[string]interface{}{
			"content": "common/jwt.go",
			"error":   "redis connect is not exists",
		})
		return
	}
	err = jwtParam.RedisCache.
		Set(context.Background(), jwtParam.TokenKey+userIdString, tokenString, jwtParam.TokenLife).
		Err()
	if err != nil {
		app_obj.GetLog().Error(c, map[string]interface{}{
			"content": "token create error",
			"error":   err.Error(),
		})
		return
	}

	return
}

func ParseToken(myToken string, ctx *base.Context) (jwtUser app_obj.JwtUser, err error) {

	logContent := make(map[string]interface{}, 6)
	logContent["token"] = "'" + myToken + "'"
	defer func() {
		if err != nil {
			logContent["error"] = err.Error()
			ctx.Error(logContent, "baseWrapperParseToken")
		}
	}()
	if myToken == "" {
		logContent["content"] = "token is null"
		ctx.Warn(logContent, "baseWrapperParseToken")
		return
	}
	jwtParam := app_obj.GetJwtParam()
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	})
	if err != nil {
		logContent["content"] = "parse token has error"
		return
	}

	if !token.Valid {
		logContent["content"] = fmt.Sprintf("token is invalid(%s)", myToken)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	var sub string
	var ok bool
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
	var res string
	if res, err = jwtParam.RedisCache.
		Get(context.Background(), jwtParam.TokenKey+sub).
		Result(); err != redis.Nil && err != nil {
		logContent["content"] = "get token from redis error"
		return
	}

	if res == "" || res != myToken {
		desc := "token is invalid"
		logContent["content"] = desc
		err = fmt.Errorf(desc)
		return
	}

	// refresh the token life time
	err = jwtParam.RedisCache.Set(context.Background(), jwtParam.TokenKey+sub, myToken, jwtParam.TokenLife).Err()
	if err != nil {
		logContent["content"] = "token create error"
		return
	}

	if err = json.Unmarshal([]byte(sub), &jwtUser); err != nil {
		logContent["content"] = "sub is error may be is not a json string"
		return
	}
	return
}

func UnsetToken(myToken string, c *gin.Context) (bool, error) {
	jwtParam := app_obj.GetJwtParam()

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	})
	if err != nil {
		app_obj.GetLog().Error(c, map[string]interface{}{
			"content": "parse token has error",
			"error":   err.Error(),
			"token":   "'" + myToken + "'",
		})
		return false, err
	}
	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		app_obj.GetLog().Error(c, map[string]interface{}{
			"content": "claims duan yan is error",
			"error":   err.Error(),
		})
		return false, errors.New("claims duan yan is error")
	}
	err = jwtParam.RedisCache.Del(context.Background(), jwtParam.TokenKey+sub).Err()
	if err != nil {
		app_obj.GetLog().Error(c, map[string]interface{}{
			"content": "unset token has error",
			"error":   err.Error(),
		})
		return false, err
	}
	return true, nil
}
