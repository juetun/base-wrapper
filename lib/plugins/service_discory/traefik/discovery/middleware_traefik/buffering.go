// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareBuffering struct {
	Buffering HttpBufferingArg `json:"buffering" yaml:"buffering,omitempty" key_value:"buffering,omitempty"`
}
type HttpBufferingArg struct {
	MaxRequestBodyBytes  int    `json:"max_request_body_bytes" yaml:"maxRequestBodyBytes,omitempty" key_value:"maxRequestBodyBytes,omitempty"`
	MemRequestBodyBytes  int    `json:"mem_request_body_bytes" yaml:"memRequestBodyBytes,omitempty" key_value:"memRequestBodyBytes,omitempty"`
	MaxResponseBodyBytes int    `json:"max_response_body_bytes" yaml:"maxResponseBodyBytes,omitempty" key_value:"maxResponseBodyBytes,omitempty"`
	MemResponseBodyBytes int    `json:"mem_response_body_bytes" yaml:"memResponseBodyBytes,omitempty" key_value:"memResponseBodyBytes,omitempty"`
	RetryExpression      string `json:"retry_expression" yaml:"retryExpression,omitempty" key_value:"retryExpression,omitempty"`
}
