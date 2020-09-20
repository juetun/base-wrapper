package response

import "github.com/gin-gonic/gin"

type detail struct {
}

func (c *detail) Success(ctx *gin.Context, model interface{}) {
	response(ctx, 0, "", model)
}

func (c *detail) Fail(ctx *gin.Context, msg string) {
	response(ctx, -1, msg, struct{}{})
}
