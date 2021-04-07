// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareForwardAuth struct {
	ForwardAuth HttpMiddlewareForwardAuthArg `json:"forward_auth" yaml:"forwardAuth,omitempty" key_value:"forwardAuth,omitempty"`
}
type HttpMiddlewareForwardAuthArg struct {
	Address                  string                          `json:"address" yaml:"address,omitempty" key_value:"address,omitempty"`
	Tls                      HttpMiddlewareForwardAuthArgTls `json:"tls" yaml:"tls,omitempty" key_value:"tls,omitempty"`
	TrustForwardHeader       bool                            `json:"trust_forward_header" yaml:"trustForwardHeader,omitempty" key_value:"trustForwardHeader,omitempty"`
	AuthResponseHeaders      []string                        `json:"auth_response_headers" yaml:"authResponseHeaders,omitempty" key_value:"authResponseHeaders,omitempty"`
	AuthResponseHeadersRegex string                          `json:"auth_response_headers_regex" yaml:"authResponseHeadersRegex,omitempty" key_value:"authResponseHeadersRegex,omitempty"`
	AuthRequestHeaders       []string                        `json:"auth_request_headers" yaml:"authRequestHeaders,omitempty" key_value:"authRequestHeaders,omitempty"`
}
type HttpMiddlewareForwardAuthArgTls struct {
	CA                 string `json:"ca" yaml:"ca,omitempty" key_value:"ca,omitempty"`
	CaOptional         bool   `json:"caOptional" yaml:"caOptional,omitempty" key_value:"caOptional,omitempty"`
	Cert               string `json:"cert" yaml:"cert,omitempty" key_value:"cert,omitempty"`
	Key                string `json:"key" yaml:"key,omitempty" key_value:"key,omitempty"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify" yaml:"insecureSkipVerify,omitempty" key_value:"insecureSkipVerify,omitempty"`
}
