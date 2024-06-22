package main

import (
	"github.com/3Xpl0it3r/kubecraft/pkg/apiserver"
	craftcert "github.com/3Xpl0it3r/kubecraft/pkg/certs"
	_ "github.com/3Xpl0it3r/kubecraft/pkg/install"
)

func main() {
	craftcert.TryCreateCACertAndKeyFiles()
	craftcert.TryCreateKubeConfig()
	craftcert.TryCreateApiServerCert()

	genericConfig := apiserver.NewConfig()
	apiserver := genericConfig.New()
	stopCh := make(chan struct{})
	apiserver.Run(stopCh)
}

func StartApiServer() {
	// simulate kube-controller-manager
	panic("unimplemented")
}

func StartControllerManager() {
	// simulate kube-controller-manager
	panic("unimplemented")
}

func StartScheduler() {
	// simulate kube-scheduler
	panic("unimplemented")
}
