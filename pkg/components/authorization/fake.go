package authorization

import (
	"context"

	"k8s.io/apiserver/pkg/authorization/authorizer"
)

// FakeAuthorizator represent fakeauthorizator
type fakeAuthorizer struct {
}

// Authorize implements authorizer.Authorizer.
func (*fakeAuthorizer) Authorize(ctx context.Context, a authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
	return authorizer.DecisionAllow, "", nil
}

func NewFakeAuthorizer() *fakeAuthorizer {
	return &fakeAuthorizer{}
}

var _ authorizer.Authorizer = new(fakeAuthorizer)
