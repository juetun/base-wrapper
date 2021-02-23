/**
* @Author:changjiang
* @Description:
* @File:micro_rpc_test
* @Version: 1.0.0
* @Date 2020/10/20 11:27 下午
 */
package test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/app/app_start"
	"github.com/juetun/base-wrapper/lib/base"
	_ "github.com/juetun/base-wrapper/lib/app/init" // 加载公共插件项
	. "github.com/juetun/base-wrapper/lib/plugins"  // 加载路由信息
	"github.com/juetun/base-wrapper/lib/plugins/rpc"
)

func TestMicroRpcGet(t *testing.T) {
	app_obj.BaseDirect, _ = filepath.Abs("../")
	app_start.NewPluginsOperate().Use(
		// PluginJwt, // 加载用户验证插件,必须放在Redis插件后
		// PluginElasticSearchV7,
		PluginShortMessage,
		PluginAppMap,
		// plugins.PluginOss,
		// plugins.PluginUser, // 用户登录,jwt等用户信息逻辑处理
	).LoadPlugins() // 加载插件动作

	appName := "api_user"
	ro := rpc.RequestOptions{
		Method:      "Get",
		AppName:     appName,
		PathVersion: "v1",
		URI:         "/permit/admin_user",
		Header:      http.Header{},
	}
	var data = base.Result{}
	err := rpc.NewHttpRpc(&ro).Send().
		GetBody().
		Bind(&data).Error
	t.Log(err.Error())
	// t.Log(data)
}
func TestMicroRpcPost(t *testing.T) {
	app_obj.BaseDirect, _ = filepath.Abs("../")
	appName := "api_user"
	ro := rpc.RequestOptions{
		Method:      "Post",
		AppName:     appName,
		PathVersion: "v1",
		URI:         "/permit/admin_user",
		Header:      http.Header{},
	}
	var data = base.Result{}
	err := rpc.NewHttpRpc(&ro).Send().
		GetBody().
		Bind(&data).Error
	t.Log(err.Error())
}
