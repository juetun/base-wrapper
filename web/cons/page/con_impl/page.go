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
	"fmt"
	"github.com/juetun/base-wrapper/lib/plugins/short_message_impl"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/render"
	"github.com/juetun/base-wrapper/lib/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	. "github.com/juetun/base-wrapper/lib/base"
	. "github.com/juetun/base-wrapper/lib/base/page_block"
	"github.com/juetun/base-wrapper/web/cons/page"
	"github.com/juetun/base-wrapper/web/srvs/srv_impl"
	"github.com/juetun/base-wrapper/web/wrapper"
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

// Websocket web socket操作
func (r *ConPageImpl) Websocket(c *gin.Context) {
	var arg wrapper.ArgWebSocket
	var err error
	var conn *websocket.Conn
	var commonParams ArgWebSocketBase
	if conn, commonParams, err = r.UpgradeWebsocket(c, &arg); err != nil {
		r.Response(c, 0, nil, err.Error())
		return
	}
	arg.ArgWebSocketBase = commonParams
	srv_impl.NewSrvWebSocketImpl().
		WebsocketSrv(conn, &arg)

	// conn.Request().Header.Set(app_obj.HttpTraceId,)

	// srv := srv_impl.NewSrvWebSocketImpl(CreateContext(&r.ControllerBase))

	// websocket_anvil.NewMessageService()
	// for {
	// 	var msg string
	// 	if err := websocket.Message.Receive(conn, &msg); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	log.Printf("recv: %v", msg)
	// 	go func() {
	// 		time.Sleep(time.Second * 1)
	// 		data := []byte(
	// 			"延迟发送" + time.Now().Format(time.RFC3339))
	// 		if _, err := conn.Write(data); err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 	}()
	// 	data := []byte(time.Now().Format(time.RFC3339))
	// 	if _, err := conn.Write(data); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// }

}

func (r *ConPageImpl) shortMessage(c *gin.Context) {
	keyList := short_message_impl.ShortMessageObj.GetChannelKey()
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
	var err error
	var res string
	res, err = (&utils.QrCodeParams{
		Width:         200,
		Content:       "http://www.google.com",
		TargetImgPath: "a.png",
	}).CreateQrCodeAsBase64Code()

	err = (&utils.QrCodeParams{
		Width:         200,
		Content:       "http://www.google.com",
		TargetImgPath: "a.png",
	}).CreateQrCodeToFile()
	var msg string
	if err != nil {
		msg = err.Error()
	}
	_ = msg
	c.Render(http.StatusOK, render.String{Format: "string", Data: []interface{}{res}})
	return
}

func (r *ConPageImpl) Main(c *gin.Context) {
	var err error
	var arg = wrapper.ArgumentDefault{}
	if err = c.BindQuery(&arg); err != nil {
		return
	}

	_ = srv_impl.NewServiceDefaultImpl(CreateContext(&r.ControllerBase, c))

	blockChild1 := NewBlock(Name("controller_main_1"), TempFile("a1.html"))

	blockChild2 := NewBlock(Name("controller_main_2"), TempFile("a2.html"))

	h := gin.H{"data": "haha"}
	block := NewBlock(Name("controller_main"), Data(gin.H{"data": "haha"}), TempFile("a.html"),
		CacheBlockOption(ExpireTime(80*time.Second)),
		ChildBock(blockChild1, blockChild2),
		RunAfter(func(block *Block) (err error) { return }))
	if h["show"], err = block.Run(); err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseHtml(c, r.MainTplFile, h)

}

func (r *ConPageImpl) MainSign(c *gin.Context) {
	var err error
	res, sign, err := NewSign().
		SignGinRequest(c)
	var msg string
	if err != nil {
		msg = err.Error()
	}
	r.Response(c, 0, res, msg+" sign:"+sign)

}
