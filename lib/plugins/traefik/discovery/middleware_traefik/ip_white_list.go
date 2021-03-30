// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareIpWhiteList struct {
	IpWhiteList HttpMiddlewareIpWhiteListArg `json:"ipWhiteList" yaml:"ipWhiteList,omitempty" key_value:"ipWhiteList,omitempty"`
}
type HttpMiddlewareIpWhiteListArg struct {
	SourceRange []string                               `json:"source_range" yaml:"sourceRange,omitempty" key_value:"sourceRange,omitempty"`
	IpStrategy  HttpMiddlewareIpWhiteListArgIpStrategy `json:"ip_strategy" yaml:"ipStrategy,omitempty" key_value:"ipStrategy,omitempty"`
}

type HttpMiddlewareIpWhiteListArgIpStrategy struct {
	Depth       int      `json:"depth" yaml:"depth,omitempty" key_value:"depth,omitempty"`
	ExcludedIPs []string `json:"excluded_ips" yaml:"excludedIPs,omitempty" key_value:"excludedIPs,omitempty"`
}
