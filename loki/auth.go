package loki

import "net/http"

type Authentication interface {
	Apply(req *http.Request)
}

type NoAuth struct{}

func (n *NoAuth) Apply(req *http.Request) {}

type BasicAuthentication struct {
	username string
	password string
}

func NewAuth(uname, pass string) Authentication {
	return &BasicAuthentication{
		username: uname,
		password: pass,
	}
}

func (b *BasicAuthentication) Apply(req *http.Request) {
	req.SetBasicAuth(b.username, b.password)
}
