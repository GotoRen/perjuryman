package internal

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// HandleTLS gets root certificate. This function returns TLS config.
func HandleTLS() (tlsConf *tls.Config, err error) {
	// Get RootCA information
	caPool := x509.NewCertPool()
	if pemCert, err := os.ReadFile(os.Getenv("ROOT_CERTIFICATE_NAME")); err != nil {
		return nil, err
	} else {
		if ok := caPool.AppendCertsFromPEM(pemCert); !ok {
			err = errors.New("The certificate is incorrect")
			return nil, err
		}
	}

	// Define TLS config structure
	tlsConf = &tls.Config{
		MinVersion: tls.VersionTLS13,
		RootCAs:    caPool,
	}

	return tlsConf, nil
}
