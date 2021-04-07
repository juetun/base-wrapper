// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewarePassTLSClientCert struct {
	PassTLSClientCert HttpMiddlewarePassTLSClientCertArg `json:"pass_tls_client_cert" yaml:"passTLSClientCert,omitempty" key_value:"passTLSClientCert,omitempty"`
}
type HttpMiddlewarePassTLSClientCertArg struct {
	Pem  bool                                   `json:"pem" yaml:"pem,omitempty" key_value:"pem,omitempty"`
	Info HttpMiddlewarePassTLSClientCertArgInfo `json:"info" yaml:"info,omitempty" key_value:"info,omitempty"`
}
type HttpMiddlewarePassTLSClientCertArgInfo struct {
	NotAfter     bool                                      `json:"not_after" yaml:"notAfter,omitempty" key_value:"notAfter,omitempty"`
	NotBefore    bool                                      `json:"not_before" yaml:"notBefore,omitempty" key_value:"notBefore,omitempty"`
	Sans         bool                                      `json:"sans" yaml:"sans,omitempty" key_value:"sans,omitempty"`
	Subject      HttpMiddlewarePassTLSClientCertArgSubject `json:"subject" yaml:"subject,omitempty" key_value:"subject,omitempty"`
	Issuer       HttpMiddlewarePassTLSClientCertArgSubject `json:"issuer" yaml:"issuer,omitempty" key_value:"issuer,omitempty"`
	SerialNumber bool                                      `json:"serial_number" yaml:"serialNumber,omitempty" key_value:"serialNumber,omitempty"`
}
type HttpMiddlewarePassTLSClientCertArgSubject struct {
	Country         bool `json:"country" yaml:"country,omitempty" key_value:"country,omitempty"`
	Province        bool `json:"province" yaml:"province,omitempty" key_value:"province,omitempty"`
	Locality        bool `json:"locality" yaml:"locality,omitempty" key_value:"locality,omitempty"`
	Organization    bool `json:"organization" yaml:"organization,omitempty" key_value:"organization,omitempty"`
	CommonName      bool `json:"commonName" yaml:"commonName,omitempty" key_value:"commonName,omitempty"`
	SerialNumber    bool `json:"serialNumber" yaml:"serialNumber,omitempty" key_value:"serialNumber,omitempty"`
	DomainComponent bool `json:"domainComponent" yaml:"domainComponent,omitempty" key_value:"domainComponent,omitempty"`
}
