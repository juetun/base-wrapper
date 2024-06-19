### 启动示例

#### 说明
1、redis与mysql使用的组件默认已注册到启动中，无需单独配置
2、框架采用插件的方式使用第三方组件(按需加载可根据需求定制自身需求组件)

```go

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	_ "github.com/juetun/base-wrapper/lib/app/init" // 加载公共插件项
	"github.com/juetun/base-wrapper/lib/authorization/model"
	. "github.com/juetun/base-wrapper/lib/plugins" // 组件目录
	_ "github.com/juetun/base-wrapper/web/router"  // 加载路由

	_ "github.com/juetun/base-wrapper/docs" //加载swagger文件
)

func main() {
    //启动需要加载的组件(按需要加载)
    app_start.NewPlugins(app_start.Authorization(&authorization)).Use(
        PluginRegistry,
        PluginClickHouse,
        PluginOss,
        PluginJwt, // 加载用户验证插件,必须放在Redis插件后
        // PluginElasticSearchV7,
        	short_message_impl.PluginShortMessage,
        PluginAppMap,
        PluginAuthorization,
        // func(arg *app_start.PluginsOperate) (err error) {
        // 	// 启动websocket
        // 	go anvil_websocket.WebsocketStart()
        // 	return
        // },
        // plugins.PluginOss,
    ).LoadPlugins() // 加载插件动作
    
    //静态文件需要注册的文件目录
    loadRouter := func(r *gin.Engine) (err error) {
        r.LoadHTMLGlob("web/views/**/*.htm")
        r.Static("/static/home", "./static/home")
        r.Static("/static/car", "./static/car")
        r.StaticFile("/jd_root.txt", "./static/jd_root.txt")
        r.StaticFile("/favicon.ico", "./static/favicon.ico")
        return
    }
    
    // 启动GIN服务
    _ = app_start.NewWebApplication().
        SetFlagMicro(true). //是否注册微服务注册中心动作
        LoadRouter(loadRouter). // 记载gin 路由配置
        Run()
    
    }
```
具体代码结构如下:
``` shell

```
