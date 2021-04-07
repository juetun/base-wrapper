// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareCompress struct {
	Compress HttpCompressArg `json:"compress" yaml:"compress,omitempty" key_value:"compress,omitempty"`
}

type HttpCompressArg struct {
	ExcludedContentTypes []string `json:"excluded_content_types" yaml:"excludedContentTypes,omitempty" key_value:"excludedContentTypes,omitempty"`
}
