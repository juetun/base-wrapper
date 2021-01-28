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
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	. "github.com/juetun/base-wrapper/lib/base"
	. "github.com/juetun/base-wrapper/lib/base/page_block"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/common/signencrypt"
	"github.com/juetun/base-wrapper/web/cons/page"
	"github.com/juetun/base-wrapper/web/srvs/srv_impl"
	"github.com/juetun/base-wrapper/web/wrapper"
	"golang.org/x/net/websocket"
)

type ConPageImpl struct {
	ControllerWeb
}

func NewConPage() (res page.ConPage) {
	p := &ConPageImpl{}
	p.ControllerWeb.Init()
	p.MainTplFile = "car_master.htm"
	return p
}

// web socket操作
func (r *ConPageImpl) Websocket(conn *websocket.Conn) {
	for {
		var msg string
		if err := websocket.Message.Receive(conn, &msg); err != nil {
			log.Println(err)
			return
		}
		log.Printf("recv: %v", msg)
		go func() {
			time.Sleep(time.Second * 1)
			data := []byte(
				"延迟发送" + time.Now().Format(time.RFC3339))
			if _, err := conn.Write(data); err != nil {
				log.Println(err)
				return
			}
		}()
		data := []byte(time.Now().Format(time.RFC3339))
		if _, err := conn.Write(data); err != nil {
			log.Println(err)
			return
		}
	}
}
func (r *ConPageImpl) shortMessage(c *gin.Context) {
	keyList := app_obj.ShortMessageObj.GetChannelKey()
	fmt.Println("当前支持的通道号有:", keyList)
	// app_obj.ShortMessageObj.SendMsg(&app_obj.MessageArgument{
	//	Mobile:   "",
	//	AreaCode: "86",
	//	Content:  "",
	// })

	var err error
	var arg wrapper.ArgumentDefault
	var result = NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(GetControllerBaseContext(&r.ControllerBase, c))
	result.Data, err = srv.Tmain(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}
func (r *ConPageImpl) Tsst(c *gin.Context) {

}

func getRequestParams(c *gin.Context) (valueMap map[string]string) {
	valueMap = make(map[string]string, len(c.Request.PostForm))
	c.Request.ParseMultipartForm(128) // 保存表单缓存的内存大小128M
	for k, v := range c.Request.Form {
		valueMap[k] = strings.Join(v, ";")
	}
	return
}

// 加密字符串
func sortParamsAndJoinData(data map[string]string, secret string) (res bytes.Buffer, err error) {
	if res, err = signencrypt.Sign().
		SignTopRequest(data, secret, signencrypt.CHARSET_UTF_8);
		err != nil {
		return
	}
	return
}

// http请求加密算法
func SignGinRequest(c *gin.Context) (signRes string, err error) {
	var secret string = "abc"
	var bt bytes.Buffer
	var encryptionCode bytes.Buffer
	bt.WriteString(c.Request.Method)
	bt.WriteString(c.Request.URL.Path)

	// 判断签名是否传递了时间
	if headerT := c.GetHeader("t"); headerT == "" {
		err = fmt.Errorf("the header must be include timestamp parameter(t)")
		return
	} else {
		var t int
		if t, err = strconv.Atoi(headerT); err != nil {
			return
		}
		// 传递的时间格式必须大于当前时间-一天
		if app_obj.App.AppEnv != common.ENV_RELEASE && int(time.Now().UnixNano()/1e6)-t > 86400000 {
			err = fmt.Errorf("the header of  parameter(t) must be more than now desc one days")
			return
		}
		bt.WriteString(headerT)
	}

	// 如果传JSON 单独处理
	if c.GetHeader("Content-Type") == "application/json" {
		bt.WriteString(secret)
		var body []byte
		if body, err = ioutil.ReadAll(c.Request.Body); err != nil {
			return
		}
		bt.Write(body)
	} else { // 如果是非JSON 传参
		// 如果不是JSON 则直接过去FORM表单参数
		if encryptionCode, err = sortParamsAndJoinData(getRequestParams(c), secret); err != nil {
			return
		}
		bt.Write(encryptionCode.Bytes())
	}

	encryptionString := strings.ToLower(bt.String())
	base64Code := base64.StdEncoding.EncodeToString([]byte(encryptionString))

	// 配置回调输出
	listenHandlerStruct := signencrypt.ListenHandlerStruct{}

	// 如果不是线上环境,可输出签名格式 (此处代码为调试 签名是否能正常使用准备)
	if app_obj.App.AppEnv != common.ENV_RELEASE && c.GetHeader("debug") != "" {
		c.Header("Sign-format", encryptionString)
		c.Header("Sign-Base64Code", base64Code)
		listenHandlerStruct = signencrypt.ListenHandlerStruct{
			MD5HMAC: func(s string) {
			},
			ByteTo16After: func(s string) {
				c.Header("Sign-ByteTo16", s)
			},
			FinishHandler: func(s string) {
				c.Header("Sign-f", s)
			},
		}

	}
	signRes = signencrypt.Sign().Encrypt(base64Code, secret, signencrypt.CHARSET_UTF_8, listenHandlerStruct)
	return
}

func (r *ConPageImpl) Main(c *gin.Context) {
	SignGinRequest(c)
	var err error
	var arg = wrapper.ArgumentDefault{}
	if err = c.BindQuery(&arg); err != nil {
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(GetControllerBaseContext(&r.ControllerBase, c))
	ctx := context.WithValue(context.TODO(), "srv", srv)
	blockChild1 := NewBlock(
		Ctx(ctx),
		Name("controller_main_1"),
		Data(gin.H{"data": "haha",}),
		TempFile("a1.html"),
	)
	blockChild2 := NewBlock(
		Ctx(ctx),
		Name("controller_main_2"),
		Data(gin.H{"data": "haha",}),
		TempFile("a2.html"),
	)

	h := gin.H{"data": "haha",}
	block := NewBlock(
		CacheBlockOption(ExpireTime(80*time.Second)),
		Ctx(ctx),
		Name("controller_main"),
		Data(h),
		TempFile("a.html"),
		ChildBock(blockChild1, blockChild2),
		RunAfter(func(block *Block) (err error) {
			return
		}),
	)

	if h["show"], err = block.Run(); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.ResponseHtml(c, r.MainTplFile, h)

}
