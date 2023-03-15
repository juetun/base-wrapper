package base

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
)

const (
	ControllerGetParamTypeBind             = "Bind"
	ControllerGetParamTypeBindJSON         = "BindJSON"
	ControllerGetParamTypeBindXML          = "BindXML"
	ControllerGetParamTypeBindQuery        = "BindQuery"
	ControllerGetParamTypeBindYAML         = "BindYAML"
	ControllerGetParamTypeBindHeader       = "BindHeader"
	ControllerGetParamTypeBindUri          = "BindUri"
	ControllerGetParamTypeShouldBind       = "ShouldBind"
	ControllerGetParamTypeShouldBindJSON   = "ShouldBindJSON"
	ControllerGetParamTypeShouldBindXML    = "ShouldBindXML"
	ControllerGetParamTypeShouldBindQuery  = "ShouldBindQuery"
	ControllerGetParamTypeShouldBindYAML   = "ShouldBindYAML"
	ControllerGetParamTypeShouldBindHeader = "ShouldBindHeader"
	ControllerGetParamTypeShouldBindUri    = "ShouldBindUri"
)

type (
	ControllerBase struct {
		Log *app_obj.AppLog
	}
	ParameterInterface interface {
		Default(ctx *Context) (err error)
	}
)

// ParametersAccept 当前参数接收
func (r *ControllerBase) ParametersAccept(ctx *Context, parameter ParameterInterface, getParamsTypes ...string) (haveError bool) {
	var err error
	var getParamsType = ControllerGetParamTypeShouldBind
	if len(getParamsTypes) > 0 {
		getParamsType = getParamsTypes[0]
	}
	switch getParamsType {
	case ControllerGetParamTypeShouldBind:
		err = ctx.GinContext.ShouldBind(parameter)
	case ControllerGetParamTypeBind:
		err = ctx.GinContext.Bind(parameter)
	case ControllerGetParamTypeBindJSON:
		err = ctx.GinContext.BindJSON(parameter)
	case ControllerGetParamTypeBindXML:
		err = ctx.GinContext.BindXML(parameter)
	case ControllerGetParamTypeBindQuery:
		err = ctx.GinContext.BindQuery(parameter)
	case ControllerGetParamTypeBindYAML:
		err = ctx.GinContext.BindYAML(parameter)
	case ControllerGetParamTypeBindHeader:
		err = ctx.GinContext.BindHeader(parameter)
	case ControllerGetParamTypeBindUri:
		err = ctx.GinContext.BindUri(parameter)
	case ControllerGetParamTypeShouldBindJSON:
		err = ctx.GinContext.ShouldBindJSON(parameter)
	case ControllerGetParamTypeShouldBindXML:
		err = ctx.GinContext.ShouldBindXML(parameter)
	case ControllerGetParamTypeShouldBindQuery:
		err = ctx.GinContext.ShouldBindQuery(parameter)
	case ControllerGetParamTypeShouldBindYAML:
		err = ctx.GinContext.ShouldBindYAML(parameter)
	case ControllerGetParamTypeShouldBindHeader:
		err = ctx.GinContext.ShouldBindHeader(parameter)
	case ControllerGetParamTypeShouldBindUri:
		err = ctx.GinContext.ShouldBindUri(parameter)
	default:
		err = fmt.Errorf("当前不支持你选择的获取参数类型(%s)", getParamsType)
		return
	}
	if err != nil {
		haveError = true
		r.ResponseParametersError(ctx.GinContext, err, ErrorParameterCode)
		return
	}
	if err = parameter.Default(ctx); err != nil {
		haveError = true
		r.ResponseParametersError(ctx.GinContext, err, ErrorParameterCode)
		return
	}
	return
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

// GetUser 当前登录用户的信息
func (r *ControllerBase) GetUser(c *gin.Context) (jwtUser JwtUser) {
	jwtUser = JwtUser{}
	v, e := c.Get(ContextUserObjectKey)
	if e {
		jwtUser = v.(JwtUser)
	}

	return jwtUser
}

// 升级websocket操作
var upGrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ArgWebSocketBase struct {
	WebsocketKey string `json:"websocket_key" form:"WebsocketKey"`
	Ip           string `json:"ip" form:"ip"`
}

// UpgradeWebsocket websocket.Handler 转 gin HandlerFunc
// argObject 必须为一个指针
func (r *ControllerBase) UpgradeWebsocket(c *gin.Context, argObject interface{}) (conn *websocket.Conn, commonParam ArgWebSocketBase, err error) {
	commonParam = ArgWebSocketBase{}
	if !c.IsWebsocket() {
		err = fmt.Errorf("now request is not a websocket")
		return
	}

	if err = c.ShouldBind(argObject); err != nil {
		log.Printf("new websocket request err : %v", err.Error())
		return
	}

	commonParam.WebsocketKey = c.Request.Header.Get(app_obj.WebSocketKey)
	commonParam.Ip = c.ClientIP()

	if conn, err = upGrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		log.Printf("new websocket request err : %v", err.Error())
		return
	}
	return
}

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"message"`
}

func (r *Result) SetCode(code int) (res *Result) {
	r.Code = code
	return r
}
func (r *Result) SetErrorMsg(err error) (res *Result) {
	if err == nil {
		return
	}
	switch err.(type) {
	case *ErrorRuntimeStruct:
		r.Code = err.(*ErrorRuntimeStruct).Code
	default:
		if r.Code == 0 {
			r.Code = -1
		}
	}
	r.Msg = err.Error()
	return r
}

func (r *Result) ToJsonByte() (res []byte) {
	if r == nil {
		return
	}
	res, _ = json.Marshal(r)
	return
}

func NewResult() *Result {
	return &Result{
		Code: 0,
		Msg:  "",
	}
}

func (r *ControllerBase) Response(c *gin.Context, code int, data interface{}, msg ...string) {
	r.setCommonHeader(c)
	message := strings.Join(msg, ",")
	if message == "" {
		message = "success"
	}
	c.JSON(http.StatusOK, Result{Code: code, Data: data, Msg: message})
}

// ResponseResult 处理正常结果集
func (r *ControllerBase) ResponseResult(c *gin.Context, result *Result) {
	r.setCommonHeader(c)
	c.JSON(http.StatusOK, result)
	return
}
func (r *ControllerBase) setCommonHeader(c *gin.Context) {
	c.Header(app_obj.HttpResponseAdministrator, app_obj.App.Administrator)
	c.Header(app_obj.HttpHeaderApp, app_obj.App.AppName)
	c.Header(app_obj.HttpHeaderVersion, app_obj.App.AppVersion)
}

// ResponseError 处理错误信息句柄
func (r *ControllerBase) ResponseError(c *gin.Context, err error, code ...int) {
	result := NewResult().
		SetErrorMsg(err)
	if result.Code == SuccessCode && len(code) > 0 {
		result.SetCode(code[0])
	}
	r.setCommonHeader(c)
	c.JSON(http.StatusOK, result)
	return
}

//注: 参数错误提示，框架调用请不要直接使用
func (r *ControllerBase) ResponseParametersError(c *gin.Context, err error, code ...int) {
	result := NewResult().
		SetErrorMsg(err)
	if result.Code == SuccessCode && len(code) > 0 {
		result.SetCode(code[0])
	}
	r.setCommonHeader(c)
	c.Writer.Write(result.ToJsonByte())
	return
}
func (r *ControllerBase) ResponseCommonHtml(c *gin.Context, code int, data gin.H, extName ...string) {
	defaultExt := "tmpl"
	if len(extName) > 0 {
		defaultExt = extName[0]
	}
	codeString := strconv.Itoa(code)
	r.setCommonHeader(c)
	c.HTML(http.StatusOK, codeString[0:1]+"xx."+defaultExt, data)
}
