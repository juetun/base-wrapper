// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareStripPrefixRegex struct {
	StripPrefixRegex HttpMiddlewareStripPrefixRegexArg `json:"stripPrefixRegex" yaml:"stripPrefixRegex,omitempty" key_value:"stripPrefixRegex,omitempty"`
}
type HttpMiddlewareStripPrefixRegexArg struct {
	Regex []string `json:"regex" yaml:"regex,omitempty" key_value:"regex,omitempty"`
}
