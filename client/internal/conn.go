package internal

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/GotoRen/perjuryman/client/internal/logger"
)

// HandleTLS establishes TLS connection with the server.
func HandleTLS() (conn *tls.Conn, err error) {
	var caCert []byte
	CAPool := x509.NewCertPool()

	if caCert, err = os.ReadFile("ca.pem"); err != nil {
		logger.LogErr("Could not load ca certificate", "error", err)
		return nil, err
	} else {
		if ok := CAPool.AppendCertsFromPEM(caCert); !ok {
			logger.LogErr("The certificate is incorrect.")
			return nil, err
		}
	}

	config := &tls.Config{
		RootCAs: CAPool,
	}

	conn, err = tls.Dial("tcp", os.Getenv("SERVER_IP")+":"+os.Getenv("TLS_PORT"), config)
	if err != nil {
		logger.LogErr("client: dial", "error", err)
	}

	return conn, nil
}
