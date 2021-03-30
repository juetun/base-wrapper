// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareHeaders struct {
	Headers HttpMiddlewareHeadersArg `json:"headers"  yaml:"headers,omitempty" key_value:"headers,omitempty"`
}
type HttpMiddlewareHeadersArg struct {
	CustomRequestHeaders              []string `json:"custom_request_headers" yaml:"customRequestHeaders,omitempty" key_value:"customRequestHeaders,omitempty"`
	CustomResponseHeaders             []string `json:"custom_response_headers" yaml:"customResponseHeaders,omitempty" key_value:"customResponseHeaders,omitempty"`
	AccessControlAllowCredentials     bool     `json:"access_control_allow_credentials" yaml:"accessControlAllowCredentials,omitempty" key_value:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders         []string `json:"access_control_allow_headers" yaml:"accessControlAllowHeaders,omitempty" key_value:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods         []string `json:"access_control_allow_methods" yaml:"accessControlAllowMethods,omitempty" key_value:"accessControlAllowMethods,omitempty"`
	AccessControlAllowOrigin          string   `json:"access_control_allow_origin" yaml:"accessControlAllowOrigin,omitempty" key_value:"accessControlAllowOrigin,omitempty"`
	AccessControlAllowOriginList      []string `json:"access_control_allow_origin_list" yaml:"accessControlAllowOriginList,omitempty" key_value:"accessControlAllowOriginList,omitempty"`
	AccessControlAllowOriginListRegex []string `json:"access_control_allow_origin_list_regex" yaml:"accessControlAllowOriginListRegex,omitempty" key_value:"accessControlAllowOriginListRegex,omitempty"`
	AccessControlExposeHeaders        []string `json:"access_control_expose_headers" yaml:"accessControlExposeHeaders,omitempty" key_value:"accessControlExposeHeaders,omitempty"`
	AccessControlMaxAge               int      `json:"access_control_max_age" yaml:"accessControlMaxAge,omitempty" key_value:"accessControlMaxAge,omitempty"`
	AddVaryHeader                     bool     `json:"add_vary_header" yaml:"addVaryHeader,omitempty" key_value:"addVaryHeader,omitempty"`
	AllowedHosts                      []string `json:"allowed_hosts" yaml:"allowedHosts,omitempty" key_value:"allowedHosts,omitempty"`
	HostsProxyHeaders                 []string `json:"hosts_proxy_headers" yaml:"hostsProxyHeaders,omitempty" key_value:"hostsProxyHeaders,omitempty"`
	SslRedirect                       bool     `json:"ssl_redirect" yaml:"sslRedirect,omitempty" key_value:"sslRedirect,omitempty"`
	SslTemporaryRedirect              bool     `json:"ssl_temporary_redirect" yaml:"sslTemporaryRedirect,omitempty" key_value:"sslTemporaryRedirect,omitempty"`
	SslHost                           string   `json:"ssl_host" yaml:"sslHost,omitempty" key_value:"sslHost,omitempty"`
	SslProxyHeaders                   []string `json:"ssl_proxy_headers" yaml:"sslProxyHeaders,omitempty" key_value:"sslProxyHeaders,omitempty"`
	SslForceHost                      bool     `json:"ssl_force_host" yaml:"sslForceHost,omitempty" key_value:"sslForceHost,omitempty"`
	StsSeconds                        int      `json:"sts_seconds" yaml:"stsSeconds,omitempty" key_value:"stsSeconds,omitempty"`
	StsIncludeSubdomains              bool     `json:"sts_include_subdomains" yaml:"stsIncludeSubdomains,omitempty" key_value:"stsIncludeSubdomains,omitempty"`
	StsPreload                        bool     `json:"sts_preload" yaml:"stsPreload,omitempty" key_value:"stsPreload,omitempty"`
	ForceSTSHeader                    bool     `json:"force_sts_header" yaml:"forceSTSHeader,omitempty" key_value:"forceSTSHeader,omitempty"`
	FrameDeny                         bool     `json:"frame_deny" yaml:"frameDeny,omitempty" key_value:"frameDeny,omitempty"`
	CustomFrameOptionsValue           string   `json:"custom_frame_options_value" yaml:"customFrameOptionsValue,omitempty" key_value:"customFrameOptionsValue,omitempty"`
	ContentTypeNosniff                bool     `json:"content_type_nosniff" yaml:"contentTypeNosniff,omitempty" key_value:"contentTypeNosniff,omitempty"`
	BrowserXssFilter                  bool     `json:"browser_xss_filter" yaml:"browserXssFilter,omitempty" key_value:"browserXssFilter,omitempty"`
	CustomBrowserXSSValue             string   `json:"custom_browser_xss_value" yaml:"customBrowserXSSValue,omitempty" key_value:"customBrowserXSSValue,omitempty"`
	ContentSecurityPolicy             string   `json:"content_security_policy" yaml:"contentSecurityPolicy,omitempty" key_value:"contentSecurityPolicy,omitempty"`
	PublicKey                         string   `json:"public_key" yaml:"publicKey,omitempty" key_value:"publicKey,omitempty"`
	ReferrerPolicy                    string   `json:"referrer_policy" yaml:"referrerPolicy,omitempty" key_value:"referrerPolicy,omitempty"`
	FeaturePolicy                     string   `json:"feature_policy" yaml:"featurePolicy,omitempty" key_value:"featurePolicy,omitempty"`
	IsDevelopment                     bool     `json:"is_development" yaml:"isDevelopment,omitempty" key_value:"isDevelopment,omitempty"`
}
