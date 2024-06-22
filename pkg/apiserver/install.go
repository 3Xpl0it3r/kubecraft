package apiserver

import (
	corestorage "github.com/3Xpl0it3r/kubecraft/pkg/registry/core"
	examplestorage "github.com/3Xpl0it3r/kubecraft/pkg/registry/example"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

func RegisterApiGroups(scheme *runtime.Scheme, parameterCodec runtime.ParameterCodec, codec serializer.CodecFactory, optsGetter genericregistry.RESTOptionsGetter) []*genericapiserver.APIGroupInfo {
	apiGroupInfos := []*genericapiserver.APIGroupInfo{}

	// register foo
	apiGroupInfos = append(apiGroupInfos, corestorage.NewRESTStorage(optsGetter))

	apiGroupInfos = append(apiGroupInfos, examplestorage.NewRESTStorage(optsGetter))

	return apiGroupInfos
}

func InstallAPIGroups(apiServer *genericapiserver.GenericAPIServer, apiGroups ...*genericapiserver.APIGroupInfo) error {

	for _, apiGroup := range apiGroups {
		if len(apiGroup.PrioritizedVersions) > 0 && len(apiGroup.PrioritizedVersions[0].Group) == 0 {
			apiServer.InstallLegacyAPIGroup("/api", apiGroup)
		} else {
			apiServer.InstallAPIGroup(apiGroup)
		}
	}
	return nil
}
