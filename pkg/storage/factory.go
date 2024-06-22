package storage

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/apiserver/pkg/storage/storagebackend/factory"
	"k8s.io/client-go/tools/cache"

	sqlitstore "github.com/3Xpl0it3r/kubecraft/pkg/storage/sqlite"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	kubestorageinterface "k8s.io/apiserver/pkg/storage"
)

func NewRawStorage() kubestorageinterface.Interface {
	return sqlitstore.NewSqliteStore()
}

// HookStorageFactoryRestOptionsFactory represent hookstoragefactoryrestoptionsfactory
type HackedStorageFactoryRestOptionsFactory struct {
	Options        options.EtcdOptions
	StorageFactory serverstorage.StorageFactory
}

func CreateRESTOptionsGetter(etcdOpts *options.EtcdOptions, factory serverstorage.StorageFactory) generic.RESTOptionsGetter {
	return &HackedStorageFactoryRestOptionsFactory{Options: *etcdOpts, StorageFactory: factory}
}

// GetRESTOptions [#TODO](should add some comments)
func (f *HackedStorageFactoryRestOptionsFactory) GetRESTOptions(resource schema.GroupResource) (generic.RESTOptions, error) {
	storageConfig, err := f.StorageFactory.NewConfig(resource)
	if err != nil {
		return generic.RESTOptions{}, fmt.Errorf("unable to find storage destination for %v, due to %v", resource, err.Error())
	}

	ret := generic.RESTOptions{
		StorageConfig:             storageConfig,
		Decorator:                 UndecoratedStorage,
		DeleteCollectionWorkers:   f.Options.DeleteCollectionWorkers,
		EnableGarbageCollection:   f.Options.EnableGarbageCollection,
		ResourcePrefix:            f.StorageFactory.ResourcePrefix(resource),
		CountMetricPollPeriod:     f.Options.StorageConfig.CountMetricPollPeriod,
		StorageObjectCountTracker: f.Options.StorageConfig.StorageObjectCountTracker,
	}

	return ret, nil
}

func UndecoratedStorage(
	config *storagebackend.ConfigForResource,
	resourcePrefix string,
	keyFunc func(obj runtime.Object) (string, error),
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object,
	getAttrsFunc storage.AttrFunc,
	trigger storage.IndexerFuncs,
	indexers *cache.Indexers) (storage.Interface, factory.DestroyFunc, error) {
	/* return NewRawStorage(config, newFunc, newListFunc, resourcePrefix) */
	return sqlitstore.NewSqliteStore(), nil, nil
}
