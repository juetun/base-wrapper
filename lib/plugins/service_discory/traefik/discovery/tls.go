// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package discovery

type TlsTraefik struct {
	Certificates []TlsTraefikCertificates     `json:"certificates" yaml:"certificates,omitempty" key_value:"certificates,omitempty"`
	Options      map[string]TlsTraefikOptions `json:"options" yaml:"options,omitempty" key_value:"options,omitempty"`
	Stores       map[string]TlsTraefikStores  `json:"stores" yaml:"stores,omitempty" key_value:"stores,omitempty"`
}
type TlsTraefikCertificates struct {
	CertFile string   `json:"cert_file" yaml:"certFile,omitempty" key_value:"certFile,omitempty"`
	KeyFile  string   `json:"key_file" yaml:"keyFile,omitempty" key_value:"keyFile,omitempty"`
	Stores   []string `json:"stores" yaml:"stores,omitempty" key_value:"stores,omitempty"`
}

type TlsTraefikOptions struct {
	MinVersion               string                      `json:"min_version" yaml:"minVersion,omitempty" key_value:"minVersion,omitempty"`
	MaxVersion               string                      `json:"max_version" yaml:"maxVersion,omitempty" key_value:"maxVersion,omitempty"`
	CipherSuites             []string                    `json:"cipher_suites" yaml:"cipherSuites,omitempty" key_value:"cipherSuites,omitempty"`
	CurvePreferences         []string                    `json:"curve_preferences" yaml:"curvePreferences,omitempty" key_value:"curvePreferences,omitempty"`
	ClientAuth               TlsTraefikOptionsClientAuth `json:"client_auth" yaml:"clientAuth,omitempty" key_value:"clientAuth,omitempty"`
	SniStrict                bool                        `json:"sni_strict" yaml:"sniStrict,omitempty" key_value:"sniStrict,omitempty"`
	PreferServerCipherSuites bool                        `json:"prefer_server_cipher_suites" yaml:"preferServerCipherSuites,omitempty" key_value:"preferServerCipherSuites,omitempty"`
}
type TlsTraefikOptionsClientAuth struct {
	CaFiles        []string `json:"ca_files" yaml:"caFiles,omitempty" key_value:"caFiles,omitempty"`
	ClientAuthType string   `json:"client_auth_type" yaml:"clientAuthType,omitempty" key_value:"clientAuthType,omitempty"`
}
type TlsTraefikStores struct {
	DefaultCertificate TlsTraefikStoresDefaultCertificate `json:"default_certificate" yaml:"defaultCertificate,omitempty" key_value:"defaultCertificate,omitempty"`
}
type TlsTraefikStoresDefaultCertificate struct {
	CertFile string `json:"cert_file" yaml:"certFile,omitempty" key_value:"certFile,omitempty"`
	KeyFile  string `json:"key_file" yaml:"keyFile,omitempty" key_value:"keyFile,omitempty"`
}
