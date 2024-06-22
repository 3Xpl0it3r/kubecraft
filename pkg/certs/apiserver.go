package certs

import (
	"fmt"
)

func TryCreateApiServerCert() {
	caCert, caKey, err := TryLoadCertAndKeyFromDisk(DefaultPkiPath, "ca")
	if err != nil {
		panic(fmt.Sprintf("failed load ca and key %w", err))
	}
	if _, _, err := TryLoadCertChainFromDisk(DefaultPkiPath, DefaultApiServerCertName); err == nil {
		return
	}
	CreateCertAndKeyFilesWithCA(caCert, caKey)
}

