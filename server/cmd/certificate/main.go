package main

import (
	"crypto/tls"
	"fmt"

	"github.com/GotoRen/perjuryman/server/exec"
	"github.com/GotoRen/perjuryman/server/internal"
	"github.com/GotoRen/perjuryman/server/internal/logger"
)

func main() {
	exec.LoadConf()

	_, err := initCert()
	if err != nil {
		logger.LogErr("Failed to get the server certificate information", "error", err)
	} else {
		fmt.Println("[INFO] Get TLS configuration information.")
	}
}

// initCert return certificate information.
func initCert() (serverTLSConf *tls.Config, err error) {
	_, _, err = internal.GenerateRootCA()
	if err != nil {
		return nil, err
	} else {
		fmt.Println("[INFO] Create RootCA")
	}

	_, err = internal.GenerateServerCert()
	if err != nil {
		return nil, err
	} else {
		fmt.Println("[INFO] Create server certificate")
	}

	return serverTLSConf, nil
}
