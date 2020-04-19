package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var MiddleWareComponent = []gin.HandlerFunc{
	Permission(),
}

// 加载权限验证Gin中间件
func LoadMiddleWare(privateMiddleWares ...gin.HandlerFunc) {
	fmt.Println("Load gin middleWare start")
	MiddleWareComponent = append(MiddleWareComponent, privateMiddleWares...)
	fmt.Println("Load gin middleWare over")

}
