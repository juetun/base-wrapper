// Package outernet
// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package outernet

import (
	"github.com/gin-gonic/gin"
)

type ConDefault interface {
	TestEs(c *gin.Context)

	Index(c *gin.Context)

	Auth(c *gin.Context)

	AuthRes(c *gin.Context)
}
