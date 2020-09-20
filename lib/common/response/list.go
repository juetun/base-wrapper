package response

import (
	"github.com/gin-gonic/gin"
)

type list struct {
}

func (c *list) Success(ctx *gin.Context, modelList *Pager) {
	response(ctx, 0, "", modelList)
}

func (c *list) Fail(ctx *gin.Context, msg string) {
	response(ctx, -1, msg, struct{}{})
}
