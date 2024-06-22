package foo

import (
	exampleapi "github.com/3Xpl0it3r/kubecraft/pkg/apis/example"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"

	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

const (
	DefaultResourceQualifiedResource = "foos"
	ResourceSigularQualifiedResource = "foo"
)


func NewStorage(optsGetter generic.RESTOptionsGetter) (*genericregistry.Store, error) {
	strategy := NewFooStrategy()
	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return new(exampleapi.Foo) },     // internal version
		NewListFunc:               func() runtime.Object { return new(exampleapi.FooList) }, // internal version
		PredicateFunc:             Matcher,
		DefaultQualifiedResource:  exampleapi.Resource(DefaultResourceQualifiedResource),
		SingularQualifiedResource: exampleapi.Resource(ResourceSigularQualifiedResource),
		CreateStrategy:            strategy,
		UpdateStrategy:            strategy,
		DeleteStrategy:            strategy,
		TableConvertor:            rest.NewDefaultTableConvertor(exampleapi.Resource(DefaultResourceQualifiedResource)),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return store, nil
}
