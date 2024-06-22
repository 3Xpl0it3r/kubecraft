package example

import (
	foostorage "github.com/3Xpl0it3r/kubecraft/pkg/registry/example/foo"
	exampleapiv1 "github.com/3Xpl0it3r/kubecraft/pkg/apis/example/v1"
	"github.com/3Xpl0it3r/kubecraft/pkg/api/legacyscheme"
	barstorage "github.com/3Xpl0it3r/kubecraft/pkg/registry/example/bar"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"

	exampleapi "github.com/3Xpl0it3r/kubecraft/pkg/apis/example"
)

func NewRESTStorage(optsGetter generic.RESTOptionsGetter) *genericapiserver.APIGroupInfo {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(exampleapi.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	storage := map[string]rest.Storage{}
	/* if store, err := foostorage.NewStorage(optsGetter); err == nil {
		storage["foos"] = store
	} */
	if store, err := barstorage.NewStorage(optsGetter); err == nil {
		storage["bars"] = store
	}
	if store, err := foostorage.NewStorage(optsGetter); err == nil {
		storage["foos"] = store
	}

	apiGroupInfo.VersionedResourcesStorageMap[exampleapiv1.SchemeGroupVersion.Version] = storage

	return &apiGroupInfo
}
