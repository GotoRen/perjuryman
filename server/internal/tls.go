package internal

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"github.com/GotoRen/perjuryman/server/internal/logger"
)

// GetCert returns serverTLSConf to get server public and private key.
func GetCert() (tlsConf *tls.Config, err error) {
	server_cert_name := "./" + os.Getenv("NAME") + ".pem"
	server_privatekey_name := "./" + os.Getenv("NAME") + ".key"

	serverCert, err := tls.LoadX509KeyPair(server_cert_name, server_privatekey_name)
	if err != nil {
		logger.LogErr("Failed to load the public and private key pair", "error", err)
		return nil, err
	}

	CAPool := x509.NewCertPool()

	if caCert, err := os.ReadFile("ca.pem"); err != nil {
		logger.LogErr("Could not load ca certificate", "error", err)
	} else {
		if ok := CAPool.AppendCertsFromPEM(caCert); !ok {
			err = errors.New("the certificate is not correct")
			logger.LogErr("The certificate is not correct.", "error", err)
			return nil, err
		}
	}

	tlsConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		MinVersion:   tls.VersionTLS13,
	}

	return tlsConf, nil
}
