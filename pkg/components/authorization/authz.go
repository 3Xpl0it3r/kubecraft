package authorization

import (
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/authorization/union"
)

func NewAuthorizer() authorizer.Authorizer {

	var authorizers []authorizer.Authorizer

	authorizers = append(authorizers, NewFakeAuthorizer())

	authorizer := union.New(authorizers...)
	return authorizer
}
