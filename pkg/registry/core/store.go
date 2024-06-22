package core

import (
	"github.com/3Xpl0it3r/kubecraft/pkg/api/legacyscheme"
	nodestorage "github.com/3Xpl0it3r/kubecraft/pkg/registry/core/node/storage"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

func NewRESTStorage(optsGetter generic.RESTOptionsGetter) *genericapiserver.APIGroupInfo {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo("", legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	storage := map[string]rest.Storage{}

	if nodeStorage, err := nodestorage.NewStorage(optsGetter); err != nil {
		klog.Errorf("failed to build node storage: %w", err)
	} else {
		storage["nodes"] = nodeStorage.Node
	}

	apiGroupInfo.VersionedResourcesStorageMap["v1"] = storage

	return &apiGroupInfo
}
