// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareRedirectScheme struct {
	RedirectScheme HttpMiddlewareRedirectSchemeArg `json:"redirect_scheme" yaml:"redirectScheme,omitempty" key_value:"redirectScheme,omitempty"`
}
type HttpMiddlewareRedirectSchemeArg struct {
	Scheme    string `json:"scheme" yaml:"scheme,omitempty" key_value:"scheme,omitempty"`
	Port      string `json:"port" yaml:"port,omitempty" key_value:"port,omitempty"`
	Permanent bool   `json:"permanent" yaml:"permanent,omitempty" key_value:"permanent,omitempty"`
}
