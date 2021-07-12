package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
)

func init() {
	app_start.HandleFuncAdminNet = append(app_start.HandleFuncAdminNet,
		func(r *gin.Engine, urlPrefix string) {

		})
}
