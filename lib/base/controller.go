package base

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
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

// GetUser 当前登录用户的信息
func (r *ControllerBase) GetUser(c *gin.Context) (jwtUser app_obj.JwtUserMessage) {
	jwtUser = app_obj.JwtUserMessage{}
	v, e := c.Get(app_obj.ContextUserObjectKey)
	if e {
		jwtUser = v.(app_obj.JwtUserMessage)
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

func NewResult() *Result {
	return &Result{
		Code: 0,
		Msg:  "",
	}
}

func (r *ControllerBase) Response(c *gin.Context, code int, data interface{}, msg ...string) {
	c.Header(app_obj.HttpResponseAdministrator, app_obj.App.Administrator)
	c.JSON(http.StatusOK, Result{Code: code, Data: data, Msg: strings.Join(msg, ",")})
}

// ResponseResult 处理正常结果集
func (r *ControllerBase) ResponseResult(c *gin.Context, result *Result) {
	c.Header(app_obj.HttpResponseAdministrator, app_obj.App.Administrator)
	c.JSON(http.StatusOK, result)
	return
}

// ResponseError 处理错误信息句柄
func (r *ControllerBase) ResponseError(c *gin.Context, err error, code ...int) {
	result := NewResult().
		SetErrorMsg(err)
	if len(code) > 0 {
		result.SetCode(code[0])
	}
	c.Header(app_obj.HttpResponseAdministrator, app_obj.App.Administrator)
	c.JSON(http.StatusOK, result)
	return
}

func (r *ControllerBase) ResponseCommonHtml(c *gin.Context, code int, data gin.H, extName ...string) {
	defaultExt := "tmpl"
	if len(extName) > 0 {
		defaultExt = extName[0]
	}
	codeString := strconv.Itoa(code)
	c.Header(app_obj.HttpResponseAdministrator, app_obj.App.Administrator)
	c.HTML(http.StatusOK, codeString[0:1]+"xx."+defaultExt, data)
}
