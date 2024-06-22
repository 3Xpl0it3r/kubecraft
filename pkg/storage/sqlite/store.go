package sqlite

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	k8sstoreinterface "k8s.io/apiserver/pkg/storage"
	"k8s.io/klog/v2"
)

// store represent store
type store struct {
}

func NewSqliteStore() *store {
	return &store{}
}

// Count implements storage.Interface.
func (*store) Count(key string) (int64, error) {
	return 0, nil
}

// Create implements storage.Interface.
func (*store) Create(ctx context.Context, key string, obj runtime.Object, out runtime.Object, ttl uint64) error {
	klog.Infof("create resource %v", key)
	return nil
}

// Delete implements storage.Interface.
func (*store) Delete(ctx context.Context, key string, out runtime.Object, preconditions *k8sstoreinterface.Preconditions, validateDeletion k8sstoreinterface.ValidateObjectFunc, cachedExistingObject runtime.Object) error {
	klog.Infof("delete resource %v", key)
	return nil
}

// Get implements storage.Interface.
func (*store) Get(ctx context.Context, key string, opts k8sstoreinterface.GetOptions, objPtr runtime.Object) error {
	klog.Infof("delete resource %v", key)
	return nil
}

// GetList implements storage.Interface.
func (*store) GetList(ctx context.Context, key string, opts k8sstoreinterface.ListOptions, listObj runtime.Object) error {
	klog.Infof("get  resource list %v", key)
	return nil
}

// GuaranteedUpdate implements storage.Interface.
func (*store) GuaranteedUpdate(ctx context.Context, key string, destination runtime.Object, ignoreNotFound bool, preconditions *k8sstoreinterface.Preconditions, tryUpdate k8sstoreinterface.UpdateFunc, cachedExistingObject runtime.Object) error {
	klog.Infof("update resource %v", key)
	return nil
}

// RequestWatchProgress implements storage.Interface.
func (*store) RequestWatchProgress(ctx context.Context) error {
	klog.Infof("resource watch progress")
	return nil
}

// Versioner implements storage.Interface.
func (*store) Versioner() k8sstoreinterface.Versioner {
	panic("unimplemented")
}

// Watch implements storage.Interface.
func (*store) Watch(ctx context.Context, key string, opts k8sstoreinterface.ListOptions) (watch.Interface, error) {
	return nil, nil
}

var _ k8sstoreinterface.Interface = new(store)
