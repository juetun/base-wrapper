/**
* @Author:changjiang
* @Description:
* @File:controller_web
* @Version: 1.0.0
* @Date 2020/4/20 9:12 上午
 */

// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package base

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/juetun/base-wrapper/lib/common/response"

	"github.com/gin-gonic/gin"
)

const CityCookieName = "city"
const MiddleCityCode = "city"

type ControllerWeb struct {
	ControllerBase
	*response.BaseResponse
	MainTplFile string
}

func (r *ControllerWeb) Init() {
	r.ControllerBase.Init()
	r.BaseResponse = response.NewBaseResponse()
	r.MainTplFile = "master.htm"
}

type BreadCrumb struct {
	Href   string `json:"href"`
	Label  string `json:"label"`
	Active string `json:"active"` // 加上样式
}

func Init() {

}

//渲染html
func (r *ControllerWeb) ResponseHtml(c *gin.Context, tpl string, data gin.H) {
	c.HTML(http.StatusOK, tpl, data)
}

// 获取详情的扩展名
func (r *ControllerWeb) GetDetailParamByKey(c *gin.Context, key string, ext ...string) string {
	extName := ".html"
	if len(ext) > 0 {
		extName = ext[0]
	}
	data := c.Params.ByName(key)
	return strings.TrimSuffix(data, extName)
}

func (r *ControllerWeb) GetCityCode(c *gin.Context) (code string) {
	v, ok := c.Get(MiddleCityCode)
	if ok {
		code = v.(string)
	}
	return
}

func (r *ControllerWeb) ParseTemplateToString(templateFile string, data gin.H) (htmlCode string, err error) {
	var tmp1 *template.Template
	tmp1 = template.New(templateFile) // 创建一个模板对象
	// fmt.Println(reflect.TypeOf(tmp1))
	tmp1, err = tmp1.ParseFiles(templateFile) // 解析模板
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)      // 实现了读写方法的可变大小的字节缓冲
	err = tmp1.Execute(buf, data) // err = tmp1.Execute(os.Stdout,name) 表示标准输出写入到控制台
	// bytes.Buffer能够从控制台获取标准输出
	if err != nil {
		return
	}
	htmlCode = buf.String()
	return
}
