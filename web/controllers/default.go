/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:04 下午
 */
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/pojos"
	"github.com/juetun/base-wrapper/web/services"
)

type ControllerDefault struct {
	base.ControllerBase
}

func NewControllerDefault() (p *ControllerDefault) {
	p = &ControllerDefault{}
	p.ControllerBase.Init()
	return p
}
func (r *ControllerDefault) Index(c *gin.Context) {
	c.Set("trace_id", "b6b200e0-2271-42fa-957d-340cfcf65f08")

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
	result.Data, err = srv.Index(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}