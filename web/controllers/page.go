/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:04 下午
 */

// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/pojos"
	"github.com/juetun/base-wrapper/web/services"
)

type ControllerPage struct {
	base.ControllerWeb
}

func NewControllerPage() (p *ControllerPage) {
	p = &ControllerPage{}
	p.ControllerBase.Init()
	return p
}
func (r *ControllerPage) Main(c *gin.Context) {
	var err error
	var arg pojos.ArgumentDefault
	var result = base.NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := services.NewServiceDefault(base.GetControllerBaseContext(&r.ControllerBase, c))
	result.Data, err = srv.TestEs(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}
