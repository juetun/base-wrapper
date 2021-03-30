// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareRetry struct {
	Retry HttpMiddlewareRetryArg `json:"retry" yaml:"retry,omitempty" key_value:"retry,omitempty"`
}
type HttpMiddlewareRetryArg struct {
	Attempts        int `json:"attempts" yaml:"attempts,omitempty" key_value:"attempts,omitempty"`
	InitialInterval int `json:"initialInterval" yaml:"initialInterval,omitempty" key_value:"initialInterval,omitempty"`
}
