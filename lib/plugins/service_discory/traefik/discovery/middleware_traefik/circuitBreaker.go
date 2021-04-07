// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareCircuitBreaker struct {
	CircuitBreaker HttpCircuitBreakerArg `json:"circuit_breaker" yaml:"circuitBreaker,omitempty" key_value:"circuitBreaker,omitempty"`
}
type HttpCircuitBreakerArg struct {
	Expression string `json:"expression" yaml:"expression,omitempty" key_value:"expression,omitempty"`
}
