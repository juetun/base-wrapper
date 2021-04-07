// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareErrors struct {
	Errors HttpMiddlewareErrorsArg `json:"errors" yaml:"errors,omitempty" key_value:"errors,omitempty"`
}
type HttpMiddlewareErrorsArg struct {
	Status  []string `json:"status" yaml:"status,omitempty" key_value:"status,omitempty"`
	Service string   `json:"service" yaml:"service,omitempty" key_value:"service,omitempty"`
	Query   string   `json:"query" yaml:"query,omitempty" key_value:"query,omitempty"`
}
