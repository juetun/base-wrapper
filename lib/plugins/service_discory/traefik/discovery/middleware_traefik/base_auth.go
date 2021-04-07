// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareBasicAuth struct {
	BasicAuth HttpMiddlewareBasicAuthArg `json:"basic_auth" yaml:"basicAuth,omitempty" key_value:"basicAuth,omitempty"`
}
type HttpMiddlewareBasicAuthArg struct {
	Users        []string `json:"users" yaml:"users,omitempty" key_value:"users,omitempty"`
	UsersFile    string   `json:"users_file" yaml:"usersFile,omitempty" key_value:"usersFile,omitempty"`
	Realm        string   `json:"realm" yaml:"realm,omitempty" key_value:"realm,omitempty"`
	RemoveHeader bool     `json:"remove_header" yaml:"removeHeader,omitempty" key_value:"removeHeader,omitempty"`
	HeaderField  string   `json:"header_field" yaml:"headerField,omitempty" key_value:"headerField,omitempty"`
}
