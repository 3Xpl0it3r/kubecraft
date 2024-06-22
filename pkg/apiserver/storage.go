package apiserver

import (
	hackedstorage "github.com/3Xpl0it3r/kubecraft/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/registry/generic"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	storagebackend "k8s.io/apiserver/pkg/storage/storagebackend"
)

var SpecialDefaultResourcePrefixes = map[schema.GroupResource]string{}

func CreateRESTOptionsGetter(Codecs *serializer.CodecFactory, Scheme *runtime.Scheme) generic.RESTOptionsGetter {
	etcdOptions := genericoptions.NewEtcdOptions(storagebackend.NewDefaultConfig("/k8sregistry", nil))
	etcdOptions.DefaultStorageMediaType = "application/vnd.kubernetes.protobuf"
	storageFactory := serverstorage.NewDefaultStorageFactory(
		etcdOptions.StorageConfig,
		etcdOptions.DefaultStorageMediaType,
		Codecs,
		defaultResourceEncodingConfigs(Scheme),
		defaultAPIResourceConfigSource(),
		SpecialDefaultResourcePrefixes)
	return hackedstorage.CreateRESTOptionsGetter(etcdOptions, storageFactory)
}

func defaultResourceEncodingConfigs(scheme *runtime.Scheme) *serverstorage.DefaultResourceEncodingConfig {
	resourceEncodingConfig := serverstorage.NewDefaultResourceEncodingConfig(scheme)
	return resourceEncodingConfig
}

func defaultAPIResourceConfigSource() *serverstorage.ResourceConfig {
	return serverstorage.NewResourceConfig()
}
