文档来源:https://ld000.space/post/traefik_add_middleware/

版本：2.2.8

Traefik 默认带很多插件，但是可能一些我们的个性化需求原生插件并不支持，这时候就需要自己开发插件了。在 2.3 版本之前 Traefik 不支持外挂插件，所以如果要添加插件的话我们需要修改源码。

下面就以添加个验证token的插件作为演示。

这个插件获取请求在header中添加的token，之后请求后端服务校验token是否正确，正确就继续请求后端，错误就直接返回错误信息。

代码修改
我们要修改3个地方，

添加插件执行文件
在 pkg/middleware/auth文件夹中添加插件主逻辑文件，这个位置可以根据自己需求修改。

![图标](https://void.oss-cn-beijing.aliyuncs.com/img/20201014142902.png)
```
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	tokenTypeName = "TokenAuthType"
)

type tokenAuth struct {
	address             string
	next                http.Handler
	name                string
	client              http.Client
}

type commonResponse struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}

// NewToken creates a passport auth middleware.
func NewToken(ctx context.Context, next http.Handler, config dynamic.TokenAuth, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, tokenTypeName)).Debug("Creating middleware")

  // 插件结构体
	ta := &tokenAuth{
		address:             config.Address,
		next:                next,
		name:                name,
	}

	// 创建请求其他服务的 http client
	ta.client = http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}

	return ta, nil
}

func (ta *tokenAuth) GetTracingInformation() (string, ext.SpanKindEnum) {
	return ta.name, ext.SpanKindRPCClientEnum
}

func (ta tokenAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), ta.name, tokenTypeName))

	errorMsg := []byte("{\"code\":10000,\"message\":\"token校验失败！\"}")

  // 从 header 中获取 token
	token := req.Header.Get("token")
	if token == "" {
		logMessage := fmt.Sprintf("Error calling %s. Cause token is empty", ta.address)
		traceAndResponseDebug(logger, rw, req, logMessage, []byte("{\"statue\":10000,\"message\":\"token is empty\"}"), http.StatusBadRequest)
		return
	}

  // 以下都是请求其他服务验证 token

	// 构建请求体
	form := url.Values{}
	form.Add("token", token)
	passportReq, err := http.NewRequest(http.MethodPost, ta.address, strings.NewReader(form.Encode()))
	tracing.LogRequest(tracing.GetSpan(req), passportReq)
	if err != nil {
		logMessage := fmt.Sprintf("Error calling %s. Cause %s", ta.address, err)
		traceAndResponseDebug(logger, rw, req, logMessage, errorMsg, http.StatusBadRequest)
		return
	}

	tracing.InjectRequestHeaders(req)

	passportReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// post 请求
	passportResponse, forwardErr := ta.client.Do(passportReq)
	if forwardErr != nil {
		logMessage := fmt.Sprintf("Error calling %s. Cause: %s", ta.address, forwardErr)
		traceAndResponseError(logger, rw, req, logMessage, errorMsg, http.StatusBadRequest)
		return
	}

	logger.Info(fmt.Sprintf("Passport auth calling %s. Response: %+v", ta.address, passportResponse))

	// 读 body
	body, readError := ioutil.ReadAll(passportResponse.Body)
	if readError != nil {
		logMessage := fmt.Sprintf("Error reading body %s. Cause: %s", ta.address, readError)
		traceAndResponseError(logger, rw, req, logMessage, errorMsg, http.StatusBadRequest)
		return
	}
	defer passportResponse.Body.Close()

	if passportResponse.StatusCode != http.StatusOK {
		logMessage := fmt.Sprintf("Remote error %s. StatusCode: %d", ta.address, passportResponse.StatusCode)
		traceAndResponseDebug(logger, rw, req, logMessage, errorMsg, http.StatusBadRequest)
		return
	}

	// 解析 body
	var commonRes commonResponse
	err = json.Unmarshal(body, &commonRes)
	if err != nil {
		logMessage := fmt.Sprintf("Body unmarshal error. Body: %s", body)
		traceAndResponseError(logger, rw, req, logMessage, errorMsg, http.StatusBadRequest)
		return
	}

	// 判断返回值，非0代表验证失败
	if commonRes.Status != 0 {
		logMessage := fmt.Sprintf("Body status is not success. Status: %d", commonRes.Status)
		traceAndResponseDebug(logger, rw, req, logMessage, errorMsg, http.StatusBadRequest)
		return
	}

	ta.next.ServeHTTP(rw, req)
}

func traceAndResponseDebug(logger log.Logger, rw http.ResponseWriter, req *http.Request, logMessage string, errorMsg []byte, status int) {
	logger.Debug(logMessage)
	tracing.SetErrorWithEvent(req, logMessage)

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.WriteHeader(status)
	_, _ = rw.Write(errorMsg)
}

func traceAndResponseInfo(logger log.Logger, rw http.ResponseWriter, req *http.Request, logMessage string, errorMsg []byte, status int) {
	logger.Info(logMessage)
	tracing.SetErrorWithEvent(req, logMessage)

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.WriteHeader(status)
	_, _ = rw.Write(errorMsg)
}

func traceAndResponseError(logger log.Logger, rw http.ResponseWriter, req *http.Request, logMessage string, errorMsg []byte, status int) {
	logger.Debug(logMessage)
	tracing.SetErrorWithEvent(req, logMessage)

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.WriteHeader(status)
	_, _ = rw.Write(errorMsg)
}
添加动态配置映射
这里添加配置文件和实体的映射关系

// pkg/config/dynamic/middlewares.go

package dynamic

/* ... */

// Middleware holds the Middleware configuration.
type Middleware struct {
  /* ... */

  // 
	TokenAuth         *TokenAuth         `json:"tokenAuth,omitempty" toml:"tokenAuth,omitempty" yaml:"tokenAuth,omitempty"`
}

/* ... */

// TokenAuth
type TokenAuth struct {
	Address             string     `json:"address,omitempty" toml:"address,omitempty" yaml:"address,omitempty"`
}
```
构造插件示例
这里写创建插件实体的代码
```
// pkg/server/middleware/middlewares.go

func (b *Builder) buildConstructor(ctx context.Context, middlewareName string) (alice.Constructor, error) {
	/* ... */

	// TokenAuth
	if config.TokenAuth != nil {
		if middleware != nil {
			return nil, badConf
		}
		middleware = func(next http.Handler) (http.Handler, error) {
			return auth.NewToken(ctx, next, *config.TokenAuth, middlewareName)
		}
	}

	/* ... */
}
```
打包配置
直接用自带的打包命令打linux包(需要安装docker)。

make binary
之后会在dist文件夹下生成可执行文件。

![图标](https://void.oss-cn-beijing.aliyuncs.com/img/20201014154636.png)

添加插件配置

http:
  middlewares:
    # token验证
    token-auth:
      tokenAuth:
        address: http://xxx.xxx.com/token_info
增加动态路由配置

http:
  routers:
    svc:
      entryPoints:
      - web
      middlewares:
      - token-auth
      service: svc
      rule: PathPrefix(`/list`)
这样新添加的插件就能用了。

插件
插件在包 pkg/middleware/XXX 里
打包命令：
script/crossbinary-default，会打多个操作系统的包，可根据实际情况注释掉对应的系统

