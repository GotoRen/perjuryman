package internal

import (
	"crypto/tls"
	"os"
)

// HandleTLS gets server certificate. This function returns TLS config.
func HandleTLS() (tlsConf *tls.Config, err error) {
	server_cert_name := os.Getenv("SERVER_CERTIFICATE_NAME")
	server_privatekey_name := os.Getenv("SERVER_CERTIFICATE_PRIVATEKEY_NAME")

	// Get server certificate
	serverCert, err := tls.LoadX509KeyPair(server_cert_name, server_privatekey_name)
	if err != nil {
		return nil, err
	}

	// Define TLS config structure
	tlsConf = &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{serverCert},
	}

	return tlsConf, nil
}
