/**
* @Author:changjiang
* @Description:
* @File:gin
* @Version: 1.0.0
* @Date 2020/3/19 11:19 下午
 */
package common

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
)

type ValidationMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type validate interface {
	Message() map[string]ValidationMessage
}
type Gin struct {
	C *gin.Context
}

func NewGin(c *gin.Context) *Gin {
	return &Gin{C: c}
}

func (r *Gin) Response(code int, data interface{}) {
	var o base.ControllerBase
	o.Response(r.C, code, data)
}

func (r *Gin) Validate(obj validate) bool {
	valid := validation.Validation{}
	b, err := valid.Valid(obj)
	if err != nil {
		app_obj.GetLog().Error(r.C,
			map[string]interface{}{
				"message": "valid error",
				"err":     err.Error(),
			})
		r.C.JSON(http.StatusOK, base.Result{Data: nil, Code: 400000000, Msg: err.Error()})
		return false
	}

	if !b {
		errorMaps := obj.Message()
		field := valid.Errors[0].Key
		if v, ok := errorMaps[field]; ok {
			r.C.JSON(http.StatusOK, base.Result{Data: v, Code: errorMaps[field].Code, Msg: errorMaps[field].Message})
			return b
		}
		r.C.JSON(http.StatusOK, base.Result{Data: nil, Code: 100000001, Msg: fmt.Sprintf("参数校验异常(%s)", field)})
		return b
	}
	return true
}
