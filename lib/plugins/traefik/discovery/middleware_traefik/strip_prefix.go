// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareStripPrefix struct {
	StripPrefix HttpMiddlewareStripPrefixArg `json:"stripPrefix" yaml:"stripPrefix,omitempty" key_value:"stripPrefix,omitempty"`
}
type HttpMiddlewareStripPrefixArg struct {
	Prefixes   []string `json:"prefixes" yaml:"prefixes,omitempty" key_value:"prefixes,omitempty"`
	ForceSlash bool     `json:"forceSlash" yaml:"forceSlash,omitempty" key_value:"forceSlash,omitempty"`
}
