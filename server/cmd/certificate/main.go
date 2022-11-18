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
		fmt.Println("[INFO] Server certificate setup is complete")
	}
}

// initCert return certificate information.
func initCert() (serverTLSConf *tls.Config, err error) {
	_, _, err = internal.GenerateRootCA()
	if err != nil {
		return nil, err
	}

	_, err = internal.GenerateServerCert()
	if err != nil {
		return nil, err
	}

	return serverTLSConf, nil
}
