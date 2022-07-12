package app_start

import (
	"context"
	"github.com/gin-gonic/gin"
)

//微服务注册逻辑
type MicroOperateInterface interface {
	//将服务信息注册入注册中心
	RegisterMicro(c *gin.Engine,cTxs ...context.Context)(ok bool,err error)

	//将服务信息从注册中心拿掉
	UnRegisterMicro()
}
