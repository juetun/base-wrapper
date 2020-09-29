package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct {
	create *create
	update *update
	detail *detail
	list   *list
}

func NewBaseResponse() *BaseResponse {
	return &BaseResponse{
		create: &create{},
		update: &update{},
		detail: &detail{},
		list:   &list{},
	}
}

func (b *BaseResponse) CreateResponse() *create {
	return b.create
}

func (b *BaseResponse) UpdateResponse() *update {
	return b.update
}

func (b *BaseResponse) DetailResponse() *detail {
	return b.detail
}

func (b *BaseResponse) ListResponse() *list {
	return b.list
}

func (b *BaseResponse) DeleteResponse() *update {
	return b.update
}

type httpResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func response(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, httpResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}



