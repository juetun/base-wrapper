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
package con_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/controllers"
	"github.com/juetun/base-wrapper/web/pojos"
	"github.com/juetun/base-wrapper/web/srv/srv_impl"
)

type ControllerDefaultImpl struct {
	base.ControllerBase
}

func NewControllerDefault() (res controllers.ControllerDefault) {
	p := &ControllerDefaultImpl{}
	p.ControllerBase.Init()
	return p
}
func (r *ControllerDefaultImpl) TestEs(c *gin.Context) {
	var err error
	var arg pojos.ArgumentDefault
	var result = base.NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(base.GetControllerBaseContext(&r.ControllerBase, c))
	result.Data, err = srv.TestEs(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}
func (r *ControllerDefaultImpl) Index(c *gin.Context) {
	var err error
	var arg pojos.ArgumentDefault
	var result = base.NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(base.GetControllerBaseContext(&r.ControllerBase, c))
	result.Data, err = srv.Index(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}