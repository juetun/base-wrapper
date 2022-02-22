package app_start

import "github.com/gin-gonic/gin"

//微服务注册逻辑
type MicroOperateInterface interface {
	//将服务信息注册入注册中心
	RegisterMicro(c *gin.Engine)(ok bool,err error)

	//将服务信息从注册中心拿掉
	UnRegisterMicro()
}
