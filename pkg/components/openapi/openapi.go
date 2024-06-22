package openapi

import (
	craftopenapi "github.com/3Xpl0it3r/kubecraft/pkg/generated/openapi"
	"k8s.io/apimachinery/pkg/runtime"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	openapicommon "k8s.io/kube-openapi/pkg/common"
)

func NewOpenApiConfig(scheme runtime.Scheme) *openapicommon.Config {

	namer := openapinamer.NewDefinitionNamer(&scheme)
	openAPIConfig := genericapiserver.DefaultOpenAPIConfig(craftopenapi.GetOpenAPIDefinitionsMegered, namer)
	return openAPIConfig
}
func NewOpenApiV3Config(scheme runtime.Scheme) *openapicommon.OpenAPIV3Config {

	namer := openapinamer.NewDefinitionNamer(&scheme)
	openAPIConfig := genericapiserver.DefaultOpenAPIV3Config(craftopenapi.GetOpenAPIDefinitionsMegered, namer)
	return openAPIConfig
}
