/**
* @Author:changjiang 当前用户所在城市
* @Description:
* @File:city
* @Version: 1.0.0
* @Date 2020/5/6 8:27 下午
 */
package plugins

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/middlewares"
)

const CityCookieName = "city"
const MiddleCityCode = "city"

func PluginLocation() (err error) {
	middlewares.MiddleWareComponent = append(middlewares.MiddleWareComponent, func(context *gin.Context) {

		// fmt.Printf("*******************")
		// fmt.Printf("读取cookie city")
		city := context.Query("city") // 优先读取URL中的城市参数
		if city == "" {
			city, _ = context.Cookie(CityCookieName)
		}
		if city == "" {
			city = "110100"
			// fmt.Printf("设置默认cookie")
		}
		config := GetCommonConfig()
		context.SetCookie(CityCookieName, city, 86400*365, "/", config.Domain, false, false)
		context.Set(MiddleCityCode, city)
	})
	return
}
