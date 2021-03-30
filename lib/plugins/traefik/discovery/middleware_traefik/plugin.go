// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewarePlugin struct {
	Plugin interface{} `json:"plugin" yaml:"plugin,omitempty" key_value:"plugin,omitempty"`
}
type HttpMiddlewarePluginArg struct {
	PluginConf HttpMiddlewarePluginConfArg `json:"plugin_conf" yaml:"PluginConf,omitempty" key_value:"PluginConf,omitempty"`
}

type HttpMiddlewarePluginConfArg struct {
	Foo string `json:"foo" yaml:"foo,omitempty" key_value:"foo,omitempty"`
}
