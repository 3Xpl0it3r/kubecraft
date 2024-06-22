//go:build tools
// +build tools

package tools

import (
	_ "k8s.io/code-generator"
	_ "k8s.io/kube-openapi/cmd/openapi-gen"
	_ "k8s.io/kube-openapi/pkg/common"
)
