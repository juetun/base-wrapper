// Package con_impl
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package con_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/cons/outernet"
	"github.com/juetun/base-wrapper/web/srvs/srv_impl"
	"github.com/juetun/base-wrapper/web/wrapper"
)

type ConDefaultImpl struct {
	base.ControllerBase
}

func NewConDefault() (res outernet.ConDefault) {
	p := &ConDefaultImpl{}
	p.ControllerBase.Init()
	return p
}

// @测试Elasticsearch
// @Description get record by ID
// @Accept  json
// @Produce json
// @Param   some_id     path    int     true        "userId"
// @Success 200 {string} string	"ok"
// @Router /record/{some_id} [get]
func (r *ConDefaultImpl) TestEs(c *gin.Context) {
	var err error
	var arg wrapper.ArgumentDefault
	var result = base.NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(base.CreateContext(&r.ControllerBase, c))
	result.Data, err = srv.TestEs(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}
func (r *ConDefaultImpl) Auth(c *gin.Context) {
	var err error
	var arg wrapper.ArgumentDefault
	var result = base.NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(base.CreateContext(&r.ControllerBase, c))
	result.Data, err = srv.Auth(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}
func (r *ConDefaultImpl) AuthRes(c *gin.Context) {
	var err error
	var arg wrapper.ArgumentDefault
	var result = base.NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(base.CreateContext(&r.ControllerBase, c))
	result.Data, err = srv.AuthRes(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}

// @获取指定ID记录
// @Description get record by ID
// @Accept  json
// @Produce json
// @Param   some_id     path    int     true        "userId"
// @Success 200 {string} string	"ok"
// @Router /record/{some_id} [get]
func (r *ConDefaultImpl) Index(c *gin.Context) {
	var err error
	var arg wrapper.ArgumentDefault
	var result = base.NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(base.CreateContext(&r.ControllerBase, c))
	result.Data, err = srv.Index(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}
