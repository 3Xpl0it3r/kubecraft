package authentication

import (
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/request/union"
)

// AuthnConfig represent authnconfig
type Authentication struct {
}

func NewAuthenticator() authenticator.Request {
	var authenticators []authenticator.Request
	// add custom audtion

	authenticators = append(authenticators, &FakeAuthenticator{})

	authenticator := union.New(authenticators...)
	return authenticator
}
