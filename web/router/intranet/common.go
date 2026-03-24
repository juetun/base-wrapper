package intranet

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	"net/http"
)

func init() {
	app_start.HandleFuncIntranet = append(app_start.HandleFuncIntranet, func(r *gin.Engine, urlPrefix string) {

		p := r.Group(urlPrefix)
		p.POST("/permit/url_path", func(c *gin.Context) {
			c.JSON(http.StatusOK, base.Result{Code: base.SuccessCode, Data: app_start.PermitAdminUrlPath, Msg: ""})
			return
		})

		p.HEAD("/heart_beat", func(c *gin.Context) {
			c.JSON(http.StatusOK, "OK")
			return
		})

	})
}
