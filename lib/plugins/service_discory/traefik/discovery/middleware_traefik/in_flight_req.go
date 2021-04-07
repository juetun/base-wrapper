// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareInFlightReq struct {
	InFlightReq HttpMiddlewareInFlightReqArg `json:"inFlightReq" yaml:"inFlightReq,omitempty" key_value:"inFlightReq,omitempty"`
}
type HttpMiddlewareInFlightReqArg struct {
	Amount          int                                         `json:"amount" yaml:"amount,omitempty" key_value:"amount,omitempty"`
	SourceCriterion HttpMiddlewareInFlightReqSourceCriterionArg `json:"source_criterion" yaml:"sourceCriterion,omitempty" key_value:"sourceCriterion,omitempty"`
}
type HttpMiddlewareInFlightReqSourceCriterionArg struct {
	IpStrategy        HttpMiddlewareInFlightReqSourceCriterionIpStrategyArg `json:"ip_strategy" yaml:"ipStrategy,omitempty" key_value:"ipStrategy,omitempty"`
	RequestHeaderName string                                                `json:"request_header_name" yaml:"requestHeaderName,omitempty" key_value:"requestHeaderName,omitempty"`
	RequestHost       bool                                                  `json:"request_host" yaml:"requestHost,omitempty" key_value:"requestHost,omitempty"`
}
type HttpMiddlewareInFlightReqSourceCriterionIpStrategyArg struct {
	Depth       int      `json:"depth" yaml:"depth,omitempty" key_value:"depth,omitempty"`
	ExcludedIPs []string `json:"excluded_ips" yaml:"excludedIPs,omitempty" key_value:"excludedIPs,omitempty"`
}
