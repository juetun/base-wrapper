// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareRedirectRegex struct {
	RedirectRegex HttpMiddlewareRedirectRegexArg `json:"redirect_regex" yaml:"redirectRegex,omitempty" key_value:"redirectRegex,omitempty"`
}
type HttpMiddlewareRedirectRegexArg struct {
	Regex       string `json:"regex" yaml:"regex,omitempty" key_value:"regex,omitempty"`
	Replacement string `json:"replacement" yaml:"replacement,omitempty" key_value:"replacement,omitempty"`
	Permanent   bool   `json:"permanent" yaml:"permanent,omitempty" key_value:"permanent,omitempty"`
}
