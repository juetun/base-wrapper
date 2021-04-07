// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareContentType struct {
	ContentType HttpContentTypeArg `json:"contentType" yaml:"contentType,omitempty" key_value:"contentType,omitempty"`
}
type HttpContentTypeArg struct {
	AutoDetect bool `json:"autoDetect" yaml:"autoDetect,omitempty" key_value:"autoDetect,omitempty"`
}
