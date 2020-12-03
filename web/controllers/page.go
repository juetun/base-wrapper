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
package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app_obj"
	. "github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/web/pojos"
	"github.com/juetun/base-wrapper/web/services"
	"golang.org/x/net/websocket"
)

type ControllerPage struct {
	ControllerWeb
}

func NewControllerPage() (p *ControllerPage) {
	p = &ControllerPage{}
	p.Init()
	p.MainTplFile = "car_master.htm"
	return
}

//web socket操作
func (r *ControllerPage) Websocket(conn *websocket.Conn) {
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
func (r *ControllerPage) shortMessage(c *gin.Context) {
	keyList := app_obj.ShortMessageObj.GetChannelKey()
	fmt.Println("当前支持的通道号有:", keyList)
	//app_obj.ShortMessageObj.SendMsg(&app_obj.MessageArgument{
	//	Mobile:   "",
	//	AreaCode: "86",
	//	Content:  "",
	//})

	var err error
	var arg pojos.ArgumentDefault
	var result = NewResult()

	err = c.ShouldBind(&arg)

	// 处理错误信息
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	srv := services.NewServiceDefault(GetControllerBaseContext(&r.ControllerBase, c))
	result.Data, err = srv.Tmain(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}

	r.ResponseResult(c, result)
}
func (r *ControllerPage) Tsst(c *gin.Context) {

}
func (r *ControllerPage) Main(c *gin.Context) {
	var err error
	var arg = pojos.ArgumentDefault{}
	if err = c.BindQuery(&arg); err != nil {
		return
	}
	blockChild1 := NewBlock(
		Name("controller_main_1"),
		Data(gin.H{"data": "haha",}),
		TempFile("a1.html"),
	)
	blockChild2 := NewBlock(
		Name("controller_main_2"),
		Data(gin.H{"data": "haha",}),
		TempFile("a2.html"),
	)

	h := gin.H{"data": "haha",}
	block := NewBlock(
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
