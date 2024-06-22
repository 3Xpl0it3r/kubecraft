package certs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/keyutil"
)

func TryCreateKubeConfig() {
	kubeConfigFilePath := filepath.Join(OutDir, DefaultKubeConfigFileName)
	if kubeConfigFileExist(kubeConfigFilePath) {
		fmt.Printf("You can connect to kubecraft via %s\n", kubeConfigFilePath)
		return
	}

	kubeConfig := CreateBasic("https://127.0.0.1:6443", DefaultCommonName, DefaultUserName, "ca")
	if err := clientcmd.WriteToFile(*kubeConfig, kubeConfigFilePath); err != nil {
		panic(err)
	}

	fmt.Printf("You can connect to kubecraft via %s\n", kubeConfigFilePath)
}

func kubeConfigFileExist(kubeConfigFile string) bool {
	_, err := os.Stat(kubeConfigFile)
	return !os.IsNotExist(err)
}

func CreateBasic(serverURL, clusterName, userName string, caName string) *clientcmdapi.Config {
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)
	caCert, caKey, err := TryLoadCertAndKeyFromDisk(DefaultPkiPath, caName)
	if err != nil {
		panic(fmt.Sprintf("failed load ca and key %w", err))
	}
	cfg := clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server:                   serverURL,
				CertificateAuthorityData: encodeCertPEM(caCert),
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		CurrentContext: contextName,
		Extensions:     map[string]runtime.Object{},
	}

	clientCert, clientKey, err := NewCertAndKey(caCert, caKey, EncryptionAlgorithmRSA2048)
	if err != nil {
		panic("failed to generate key and cert")
	}
	encodedClientKey, err := keyutil.MarshalPrivateKeyToPEM(clientKey)
	if err != nil {
		panic("faile unmarshal private key")
	}
	cfg.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		ClientKeyData:         encodedClientKey,
		ClientCertificateData: encodeCertPEM(clientCert),
	}

	return &cfg
}

func createKubeConfigFileIfNotExists(outDir, filename string, config *clientcmdapi.Config) error {
	kubeConfigFilePath := filepath.Join(outDir, filename)
	fmt.Printf("[kubeconfig] Writing %q kubeconfig file\n", filename)
	if err := clientcmd.WriteToFile(*config, kubeConfigFilePath); err != nil {

		return errors.Wrapf(err, "failed to save kubeconfig file %q on disk", kubeConfigFilePath)
	}

	return nil
}
