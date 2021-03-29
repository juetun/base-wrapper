// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareAddPrefix struct {
	AddPrefix HttpAddPrefixArg `json:"addPrefix" yaml:"addPrefix,omitempty" key_value:"addPrefix,omitempty"`
}
type HttpAddPrefixArg struct {
	Prefix string `json:"prefix" yaml:"prefix,omitempty" key_value:"prefix,omitempty"`
}
