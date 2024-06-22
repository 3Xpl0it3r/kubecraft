package authentication

import (
	"net/http"

	k8sauthenticator "k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
)

// FakeAuthenticator represent fakeauthenticator
type FakeAuthenticator struct {
}

// AuthenticateRequest implements authenticator.Request.
func (*FakeAuthenticator) AuthenticateRequest(req *http.Request) (*k8sauthenticator.Response, bool, error) {
	return &k8sauthenticator.Response{
		Audiences: []string{},
		User:      &user.DefaultInfo{},
	}, true, nil
}

var _ k8sauthenticator.Request = new(FakeAuthenticator)
