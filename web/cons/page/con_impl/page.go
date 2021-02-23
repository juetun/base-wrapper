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
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	. "github.com/juetun/base-wrapper/lib/base"
	. "github.com/juetun/base-wrapper/lib/base/page_block"
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
	srv := srv_impl.NewServiceDefaultImpl(CreateContext(&r.ControllerBase, c))
	result.Data, err = srv.Tmain(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}

func (r *ConPageImpl) Tsst(c *gin.Context) {
	_ = c
	return
}

func (r *ConPageImpl) Main(c *gin.Context) {
	var err error
	var arg = wrapper.ArgumentDefault{}
	if err = c.BindQuery(&arg); err != nil {
		return
	}
	srv := srv_impl.NewServiceDefaultImpl(CreateContext(&r.ControllerBase, c))
	ctx := context.WithValue(context.TODO(), "srv", srv)
	blockChild1 := NewBlock(
		Ctx(ctx),
		Name("controller_main_1"),
		Data(gin.H{"data": "haha"}),
		TempFile("a1.html"),
	)
	blockChild2 := NewBlock(
		Ctx(ctx),
		Name("controller_main_2"),
		Data(gin.H{"data": "haha"}),
		TempFile("a2.html"),
	)

	h := gin.H{"data": "haha"}
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
func (r *ConPageImpl) MainSign(c *gin.Context) {
	var err error
	res, sign, err := signencrypt.NewSign().
		SignGinRequest(c, func(appName string) (secret string, err error) {
			secret = "signxxx"
			// TODO通过appName获取签名值
			return
		})
	var msg string
	if err != nil {
		msg = err.Error()
	}
	r.Response(c, 0, res, msg+" sign:"+sign)

}
