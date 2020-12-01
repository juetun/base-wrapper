package base

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_obj"
)

type ControllerBase struct {
	Log *app_obj.AppLog
}

func (r *ControllerBase) Init() *ControllerBase {
	r.Log = app_obj.GetLog()
	return r
}

func (r *ControllerBase) GetOperateUser(c *gin.Context) string {
	return r.GetUser(c).Name
}
func (r *ControllerBase) GetAdminUserName(c *gin.Context) string {
	return r.GetUser(c).Name
}

// 当前登录用户的信息
func (r *ControllerBase) GetUser(c *gin.Context) (jwtUser app_obj.JwtUserMessage) {
	jwtUser = app_obj.JwtUserMessage{}
	v, e := c.Get(app_obj.ContextUserObjectKey)
	if e {
		jwtUser = v.(app_obj.JwtUserMessage)
	}

	return jwtUser
}

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"message"`
}

func (r *Result) SetErrorMsg(err error) (res *Result) {
	r.Code = -1
	r.Msg = err.Error()
	return r
}

func NewResult() *Result {
	return &Result{
		Code: 0,
		Msg:  "",
	}
}

func (r *ControllerBase) Response(c *gin.Context, code int, data interface{}, msg ...string) {
	c.JSON(http.StatusOK, Result{Code: code, Data: data, Msg: strings.Join(msg, ",")})
}

// 处理正常结果集
func (r *ControllerBase) ResponseResult(c *gin.Context, result *Result) {
	c.JSON(http.StatusOK, result)
	return
}

// 处理错误信息句柄
func (r *ControllerBase) ResponseError(c *gin.Context, err error) {
	result := NewResult().
		SetErrorMsg(err)
	c.JSON(http.StatusOK, result)
	return
}

func (r *ControllerBase) ResponseCommonHtml(c *gin.Context, code int, data gin.H, extName ...string) {
	defaultExt := "tmpl"
	if len(extName) > 0 {
		defaultExt = extName[0]
	}
	codeString := strconv.Itoa(code)
	c.HTML(http.StatusOK, codeString[0:1]+"xx."+defaultExt, data)
}
