// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareRateLimit struct {
	RateLimit HttpMiddlewareRateLimitArg `json:"rateLimit" yaml:"rateLimit,omitempty" key_value:"rateLimit,omitempty"`
}
type HttpMiddlewareRateLimitArg struct {
	Average         int                                       `json:"average" yaml:"average,omitempty" key_value:"average,omitempty"`
	Period          int                                       `json:"period" yaml:"period,omitempty" key_value:"period,omitempty"`
	Burst           int                                       `json:"burst" yaml:"burst,omitempty" key_value:"burst,omitempty"`
	SourceCriterion HttpMiddlewareRateLimitSourceCriterionArg `json:"source_criterion" yaml:"sourceCriterion,omitempty" key_value:"sourceCriterion,omitempty"`
}
type HttpMiddlewareRateLimitSourceCriterionArg struct {
	IpStrategy        HttpMiddlewareRateLimitSourceCriterionIpStrategy `json:"ip_strategy" yaml:"ipStrategy,omitempty" key_value:"ipStrategy,omitempty"`
	RequestHeaderName string                                           `json:"request_header_name" yaml:"requestHeaderName,omitempty" key_value:"requestHeaderName,omitempty"`
	RequestHost       bool                                             `json:"request_host" yaml:"requestHost,omitempty" key_value:"requestHost,omitempty"`
}
type HttpMiddlewareRateLimitSourceCriterionIpStrategy struct {
	Depth       int      `json:"depth" yaml:"depth,omitempty" key_value:"depth,omitempty"`
	ExcludedIPs []string `json:"excluded_ips" yaml:"excludedIPs,omitempty" key_value:"excludedIPs,omitempty"`
}
