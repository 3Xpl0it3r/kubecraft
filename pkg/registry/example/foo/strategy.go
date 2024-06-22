package foo

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	_ rest.RESTCreateStrategy = new(fooStrategy)
	_ rest.RESTDeleteStrategy = new(fooStrategy)
	_ rest.RESTUpdateStrategy = new(fooStrategy)
)

// fooStrategy represent foostrategy
type fooStrategy struct {
}

// AllowCreateOnUpdate implements rest.RESTUpdateStrategy.
func (e *fooStrategy) AllowCreateOnUpdate() bool {
	return true
}

// AllowUnconditionalUpdate implements rest.RESTUpdateStrategy.
func (e *fooStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// PrepareForUpdate implements rest.RESTUpdateStrategy.
func (e *fooStrategy) PrepareForUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) {
	return
}

// ValidateUpdate implements rest.RESTUpdateStrategy.
func (e *fooStrategy) ValidateUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate implements rest.RESTUpdateStrategy.
func (e *fooStrategy) WarningsOnUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) []string {
	return nil
}

func NewFooStrategy() *fooStrategy {
	return &fooStrategy{}
}

// Canonicalize implements rest.RESTCreateStrategy.
func (e fooStrategy) Canonicalize(obj runtime.Object) {
}

// GenerateName implements rest.RESTCreateStrategy.
func (e fooStrategy) GenerateName(base string) string {
	return "foos"
}

// NamespaceScoped implements rest.RESTCreateStrategy.
func (e fooStrategy) NamespaceScoped() bool {
	return true
}

// ObjectKinds implements rest.RESTCreateStrategy.
func (e fooStrategy) ObjectKinds(runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	return nil, true, nil
}

// PrepareForCreate implements rest.RESTCreateStrategy.
func (e fooStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

// Recognizes implements rest.RESTCreateStrategy.
func (e fooStrategy) Recognizes(gvk schema.GroupVersionKind) bool {
	return true
}

// Validate implements rest.RESTCreateStrategy.
func (e fooStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate implements rest.RESTCreateStrategy.
func (e fooStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}
