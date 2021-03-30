// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik
type HttpMiddlewareReplacePath struct {
	ReplacePath HttpMiddlewareReplacePathArg `json:"replacePath" yaml:"replacePath,omitempty" key_value:"replacePath,omitempty"`
}
type HttpMiddlewareReplacePathArg struct {
	Path string `json:"path" yaml:"path,omitempty" key_value:"path,omitempty"`
}
