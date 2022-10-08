一、实现一个gin框架中间件方法

```go

// SignHttp 接口签名验证
func SignHttp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var res bool
		var err error
		if res, _, err = signencrypt.NewSign().
			SignGinRequest(c, func(appName string) (secret string, err error) {
				secret = "signxxx"
				// TODO 通过appName获取签名值
				return
			}); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  "sign err",
			})
			return
		}
		if !res {
			c.AbortWithStatusJSON(http.StatusOK, base.Result{
				Code: http.StatusUnauthorized,
				Msg:  "sign validate failure",
			})
		}
		c.Next()
	}
}

```

二、 在启动入口(main.go main函数)入口调用 SignHttp方法

```go

	// 启动GIN服务
	app_start.NewWebApplication(
		SignHttp(), // 添加签名中间件
	).LoadRouter(). // 记载gin 路由配置
			Run()
```
