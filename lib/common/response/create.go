package response

import "github.com/gin-gonic/gin"

type create struct {
}

func (c *create) Success(ctx *gin.Context, hid string) {
	response(ctx, 0, "", hid)
}

func (c *create) Fail(ctx *gin.Context, msg string) {
	response(ctx, -1, msg, "")
}
