package bar

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	_ rest.RESTCreateStrategy = new(barStrategy)
	_ rest.RESTDeleteStrategy = new(barStrategy)
	_ rest.RESTUpdateStrategy = new(barStrategy)
)

// barStrategy represent barstrategy
type barStrategy struct {
}

// AllowCreateOnUpdate implements rest.RESTUpdateStrategy.
func (e *barStrategy) AllowCreateOnUpdate() bool {
	return true
}

// AllowUnconditionalUpdate implements rest.RESTUpdateStrategy.
func (e *barStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// PrepareForUpdate implements rest.RESTUpdateStrategy.
func (e *barStrategy) PrepareForUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) {
	return
}

// ValidateUpdate implements rest.RESTUpdateStrategy.
func (e *barStrategy) ValidateUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate implements rest.RESTUpdateStrategy.
func (e *barStrategy) WarningsOnUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) []string {
	return nil
}

func NewbarStrategy() *barStrategy {
	return &barStrategy{}
}

// Canonicalize implements rest.RESTCreateStrategy.
func (e barStrategy) Canonicalize(obj runtime.Object) {
}

// GenerateName implements rest.RESTCreateStrategy.
func (e barStrategy) GenerateName(base string) string {
	return "bars"
}

// NamespaceScoped implements rest.RESTCreateStrategy.
func (e barStrategy) NamespaceScoped() bool {
	return true
}

// ObjectKinds implements rest.RESTCreateStrategy.
func (e barStrategy) ObjectKinds(runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	return nil, true, nil
}

// PrepareForCreate implements rest.RESTCreateStrategy.
func (e barStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

// Recognizes implements rest.RESTCreateStrategy.
func (e barStrategy) Recognizes(gvk schema.GroupVersionKind) bool {
	return true
}

// Validate implements rest.RESTCreateStrategy.
func (e barStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate implements rest.RESTCreateStrategy.
func (e barStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}
