// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareChain struct {
	Chain HttpChainArg `json:"chain" yaml:"chain,omitempty" key_value:"chain,omitempty"`
}
type HttpChainArg struct {
	Middlewares []string `json:"middlewares" yaml:"middlewares,omitempty" key_value:"middlewares,omitempty"`
}
