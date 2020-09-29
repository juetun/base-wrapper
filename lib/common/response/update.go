package response

import "github.com/gin-gonic/gin"

type update struct {
}

func (c *update) Success(ctx *gin.Context) {
	response(ctx, 0, "", true)
}

func (c *update) Fail(ctx *gin.Context, msg string) {
	response(ctx, -1, msg, false)
}
