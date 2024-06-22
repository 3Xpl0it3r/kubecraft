package apiserver

import (
	"fmt"

	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/3Xpl0it3r/kubecraft/pkg/api/legacyscheme"
)

// ApiServer represent server
type ApiServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

func NewServer(cc *ControlplaneConfig, delegationTarget genericapiserver.DelegationTarget) *ApiServer {
	return newServer(cc, delegationTarget)
}

func newServer(cc *ControlplaneConfig, delegationTarget genericapiserver.DelegationTarget) *ApiServer {
	if delegationTarget == nil {
		delegationTarget = genericapiserver.NewEmptyDelegate()
	}
	genericApiServer, err := cc.Generic.New("generic", delegationTarget)
	if err != nil {
		panic(fmt.Sprintf("failed to create apisever %w", err))
	}

	// register apis
	apiGroups := RegisterApiGroups(legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs, cc.Generic.RESTOptionsGetter)
	InstallAPIGroups(genericApiServer, apiGroups...)

	return &ApiServer{GenericAPIServer: genericApiServer}
}

// PreRun [#TODO](should add some comments)
func (s *ApiServer) Run(stopCh chan struct{}) {
	s.GenericAPIServer.PrepareRun().Run(stopCh)
}
