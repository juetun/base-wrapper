// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package middleware_traefik

type HttpMiddlewareDigestAuth struct {
	DigestAuth HttpDigestAuthArg `json:"digestAuth" yaml:"digestAuth,omitempty" key_value:"digestAuth,omitempty"`
}
type HttpDigestAuthArg struct {
	Users        []string `json:"users" yaml:"users,omitempty" key_value:"users,omitempty"`
	UsersFile    string   `json:"users_file" yaml:"usersFile,omitempty" key_value:"usersFile,omitempty"`
	RemoveHeader bool     `json:"removeHeader" yaml:"removeHeader,omitempty" key_value:"removeHeader,omitempty"`
	Realm        string   `json:"realm" yaml:"realm,omitempty" key_value:"realm,omitempty"`
	HeaderField  string   `json:"headerField" yaml:"headerField,omitempty" key_value:"headerField,omitempty"`
}
