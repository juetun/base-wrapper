/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-14
 * Time: 23:48
 */
package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/juetun/base-wrapper/lib/app_log"
	"github.com/juetun/base-wrapper/lib/app_obj"
	utils2 "github.com/juetun/base-wrapper/lib/utils"
)

func CreateToken(user app_obj.JwtUserMessage) (tokenString string, err error) {
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
		app_log.GetLog().Error(map[string]string{
			"content": "token create error",
			"error":   err.Error(),
		})
		return
	}
	if jwtParam.RedisCache == nil {
		app_log.GetLog().Error(map[string]string{
			"content": "common/jwt.go",
			"error":   "redis connect is not exists",
		})
		return
	}
	err = jwtParam.RedisCache.
		Set(jwtParam.TokenKey+userIdString, tokenString, jwtParam.TokenLife).
		Err()
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "token create error",
			"error":   err.Error(),
		})
		return
	}

	return
}

func ParseToken(myToken string) (jwtUser app_obj.JwtUserMessage, err error) {
	jwtParam := app_obj.GetJwtParam()

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	})
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "parse token has error",
			"error":   err.Error(),
			"token":   "'" + myToken + "'",
		})
		return
	}

	if !token.Valid {
		app_log.GetLog().Error(map[string]string{
			"content": fmt.Sprintf("token is invalid(%s)", myToken),
			"error":   err.Error(),
		})
		return
	}
	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		app_log.GetLog().Error(map[string]string{
			"content": "claims duan yan is error",
			"error":   err.Error(),
		})
		err = fmt.Errorf("claims duan yan is error")
		return
	}
	if jwtParam.RedisCache == nil {
		msg := "Redis connect is null"
		app_log.GetLog().Error(map[string]string{
			"content": "common/jwt.go",
			"error":   msg,
		})
		err = fmt.Errorf(msg)
		return
	}
	res, err := jwtParam.RedisCache.
		Get(jwtParam.TokenKey + sub).
		Result()

	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "get token from redis error",
			"error":   err.Error(),
		})
		return
	}

	if res == "" || res != myToken {
		app_log.GetLog().Error(map[string]string{
			"content": "token is invalid",
			"error":   myToken,
		})
		err = fmt.Errorf("token is invalid")
		return
	}

	// refresh the token life time
	err = jwtParam.RedisCache.Set(jwtParam.TokenKey+sub, myToken, jwtParam.TokenLife).Err()
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "token create error",
			"error":   err.Error(),
		})
		return
	}
	err = json.Unmarshal([]byte(sub), &jwtUser)
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "sub is error may be is not a json string",
			"error":   err.Error(),
		})
	}
	return
}

func UnsetToken(myToken string) (bool, error) {
	jwtParam := app_obj.GetJwtParam()

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	})
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "parse token has error",
			"error":   err.Error(),
			"token":   "'" + myToken + "'",
		})
		return false, err
	}
	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		app_log.GetLog().Error(map[string]string{
			"content": "claims duan yan is error",
			"error":   err.Error(),
		})
		return false, errors.New("claims duan yan is error")
	}
	err = jwtParam.RedisCache.Del(jwtParam.TokenKey + sub).Err()
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "unset token has error",
			"error":   err.Error(),
		})
		return false, err
	}
	return true, nil
}
