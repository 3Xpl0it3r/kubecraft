package apiserver

import (
	"fmt"
	"net"
	"time"

	craftcerts "github.com/3Xpl0it3r/kubecraft/pkg/certs"

	"github.com/3Xpl0it3r/kubecraft/pkg/components/openapi"

	craftscheme "github.com/3Xpl0it3r/kubecraft/pkg/api/legacyscheme"

	"github.com/google/uuid"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	genericapiserveropts "k8s.io/apiserver/pkg/server/options"
	clientgoinformers "k8s.io/client-go/informers"
	clientgoclientset "k8s.io/client-go/kubernetes"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/component-base/version"

	"github.com/3Xpl0it3r/kubecraft/pkg/components/authentication"
	"github.com/3Xpl0it3r/kubecraft/pkg/components/authorization"
)

// Extra represent extra
type Extra struct {
}

// ControlplaneConfig represent apiserverconfig
type ControlplaneConfig struct {
	Generic *genericapiserver.CompletedConfig
}

func NewConfig() *ControlplaneConfig {
	genericConfig := genericapiserver.NewConfig(craftscheme.Codecs)

	// authenticator && authorization
	genericConfig.Authentication.Authenticator = authentication.NewAuthenticator()
	genericConfig.Authorization.Authorizer = authorization.NewAuthorizer()
	// openapi
	genericConfig.OpenAPIConfig = openapi.NewOpenApiConfig(*craftscheme.Scheme)
	genericConfig.OpenAPIV3Config = openapi.NewOpenApiV3Config(*craftscheme.Scheme)

	//secure
	key := craftcerts.ApiServerKeyFile
	cert := craftcerts.ApiServerCertFile
	ApplySecureServTo(cert, key, genericConfig)

	clientgoExternalClient, err := clientgoclientset.NewForConfig(genericConfig.LoopbackClientConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create external clientset %w", err))
	}
	versionedInformer := clientgoinformers.NewSharedInformerFactory(clientgoExternalClient, 10*time.Minute)

	kubeVersion := version.Get()
	genericConfig.Version = &kubeVersion

	// storage
	genericConfig.RESTOptionsGetter = CreateRESTOptionsGetter(&craftscheme.Codecs, craftscheme.Scheme)

	completeConfig := genericConfig.Complete(versionedInformer)

	return &ControlplaneConfig{Generic: &completeConfig}

}

func ApplySecureServTo(certFile, keyFile string, genericConfig *genericapiserver.Config) {
	listener, _, err := genericapiserveropts.CreateListener("tcp4", "127.0.0.1:6443", net.ListenConfig{})
	if err != nil {
		panic("bind addr failed")
	}

	secureConfig := genericapiserver.SecureServingInfo{}
	// 配置APIserver作为https证书
	secureConfig.Listener = listener
	secureConfig.Cert, err = dynamiccertificates.NewDynamicServingContentFromFiles("serving-cert", certFile, keyFile)
	if err != nil {
		panic("server cert is invalid")
	}

	// 生成一个自签证书用来内部通信
	certPem, keyPem, err := certutil.GenerateSelfSignedCertKey(LoopbackClientServerNameOverride, nil, nil)
	if err != nil {
		panic(fmt.Sprint("failed to create self certificate for loopbackclient"))
	}

	// 下面如果不用于实现CRD以及aggregator APIserver可以不用写进去
	// 将用于loopback的cert添加到SNI里面
	certProvider, err := dynamiccertificates.NewStaticSNICertKeyContent(LoopbackClientServerNameOverride, certPem, keyPem, LoopbackClientServerNameOverride)
	if err != nil {
		panic(fmt.Sprint("failed to create  caprovider for loopbackclient"))
	}
	secureConfig.SNICerts = append(secureConfig.SNICerts, certProvider)

	secureLookbackClient, err := secureConfig.NewLoopbackClientConfig(uuid.New().String(), certPem)
	if err != nil {
		panic(fmt.Sprint("failed to create    loopbackclient"))
	}

	genericConfig.SecureServing = &secureConfig
	genericConfig.LoopbackClientConfig = secureLookbackClient

}

// New [#TODO](should add some comments)
func (c *ControlplaneConfig) New() *ApiServer {
	return newServer(c, nil)
}
