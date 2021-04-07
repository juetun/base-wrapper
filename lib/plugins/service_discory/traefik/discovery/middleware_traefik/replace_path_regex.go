// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareReplacePathRegex struct {
	ReplacePathRegex HttpMiddlewareReplacePathRegexArg `json:"replace_path_regex" yaml:"replacePathRegex,omitempty" key_value:"replacePathRegex,omitempty"`
}
type HttpMiddlewareReplacePathRegexArg struct {
	Regex       string `json:"regex" yaml:"regex,omitempty" key_value:"regex,omitempty"`
	Replacement string `json:"replacement" yaml:"replacement,omitempty" key_value:"replacement,omitempty"`
}
