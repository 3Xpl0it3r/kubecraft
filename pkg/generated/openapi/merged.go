package openapi

import (
	k8sopenapi "github.com/3Xpl0it3r/kubecraft/pkg/generated/openapi/k8s"
	"k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitionsMegered(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {

	openapi := k8sopenapi.GetOpenAPIDefinitions(ref)
	for key, define := range GetOpenAPIDefinitions(ref) {
		openapi[key] = define
	}
	return openapi
}
