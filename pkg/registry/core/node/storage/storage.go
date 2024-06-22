/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"context"
	"net/http"
	"net/url"

	"github.com/3Xpl0it3r/kubecraft/pkg/registry/core/node"
	noderest "github.com/3Xpl0it3r/kubecraft/pkg/registry/core/node/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	api "k8s.io/kubernetes/pkg/apis/core"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// NodeStorage includes storage for nodes and all sub resources.
type NodeStorage struct {
	Node   *REST
	Status *StatusREST
	Proxy  *noderest.ProxyREST
}

// REST implements a RESTStorage for nodes.
type REST struct {
	*genericregistry.Store
	proxyTransport http.RoundTripper
}

// StatusREST implements the REST endpoint for changing the status of a node.
type StatusREST struct {
	store *genericregistry.Store
}

// New creates a new Node object.
func (r *StatusREST) New() runtime.Object {
	return &api.Node{}
}

// Destroy cleans up resources on shutdown.
func (r *StatusREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// GetResetFields implements rest.ResetFieldsStrategy
func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.store.GetResetFields()
}

func (r *StatusREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.store.ConvertToTable(ctx, object, tableOptions)
}

// NewStorage returns a NodeStorage object that will work against nodes.
func NewStorage(optsGetter generic.RESTOptionsGetter) (*NodeStorage, error) {
	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return &api.Node{} },
		NewListFunc:               func() runtime.Object { return &api.NodeList{} },
		PredicateFunc:             node.MatchNode,
		DefaultQualifiedResource:  api.Resource("nodes"),
		SingularQualifiedResource: api.Resource("node"),

		CreateStrategy:      node.Strategy,
		UpdateStrategy:      node.Strategy,
		DeleteStrategy:      node.Strategy,
		ResetFieldsStrategy: node.Strategy,

		TableConvertor: rest.NewDefaultTableConvertor(api.Resource("nodes")),
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    node.GetAttrs,
	}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = node.StatusStrategy
	statusStore.ResetFieldsStrategy = node.StatusStrategy

	// Set up REST handlers
	nodeREST := &REST{Store: store}
	statusREST := &StatusREST{store: &statusStore}
	proxyREST := &noderest.ProxyREST{Store: store}

	return &NodeStorage{
		Node:   nodeREST,
		Status: statusREST,
		Proxy:  proxyREST,
	}, nil
}

// Implement Redirector.
var _ = rest.Redirector(&REST{})

// ResourceLocation returns a URL to which one can send traffic for the specified node.
func (r *REST) ResourceLocation(ctx context.Context, id string) (*url.URL, http.RoundTripper, error) {
	return nil, nil, nil
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"no"}
}
