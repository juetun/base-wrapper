package app_obj

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

const TmpSignKey = "jueTungygoaesctr"

func GetHeaderAppName(c *gin.Context) (appName, secret string, err error) {
	URI := strings.TrimPrefix(c.Request.URL.Path, "/")
	if URI == "" {
		err = fmt.Errorf("get app name failure")
		return
	}
	urlString := strings.Split(URI, "/")
	appName = urlString[0]
	secret = TmpSignKey
	// TODO 通过appName获取签名值
	return
}
