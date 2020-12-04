/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:04 下午
 */
package cons

import (
	"github.com/gin-gonic/gin"
)

type ConDefault interface {
	TestEs(c *gin.Context)
	Index(c *gin.Context)
}
