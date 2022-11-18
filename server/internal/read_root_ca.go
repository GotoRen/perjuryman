package internal

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/GotoRen/perjuryman/server/internal/logger"
)

// ==================================================================//
// Read "RootCA" and "Root certificate"
// ==================================================================//

// GetRootCA gets root certificate and its private key.
func GetRootCA() (rootCert *x509.Certificate, rootCertPrivKey *rsa.PrivateKey, err error) {
	rootCert, err = readRootCertificate()
	if err != nil {
		return nil, nil, err
	}

	rootCertPrivKey, err = readRootCertificatePrivateKey()
	if err != nil {
		return nil, nil, err
	}

	return rootCert, rootCertPrivKey, nil
}

// readRootCertificate reads root certificate.
func readRootCertificate() (rootCert *x509.Certificate, err error) {
	f, err := os.ReadFile(os.Getenv("ROOT_CERTIFICATE_NAME"))
	if err != nil {
		logger.LogErr("RootCA: Could not load ca certificate", "error", err)
		return nil, err
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(f); !ok {
		return nil, err
	}

	certBlock, _ := pem.Decode(f)
	if certBlock == nil {
		err = errors.New("invalid private key data")
		return nil, err
	}

	rootCert, err = x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return rootCert, nil
}

// readRootCertificatePrivateKey reads root certificate's RSA private key.
func readRootCertificatePrivateKey() (rootCertPrivKey *rsa.PrivateKey, err error) {
	f, err := os.ReadFile(os.Getenv("ROOT_CERTIFICATE_PRIVATEKEY_NAME"))
	if err != nil {
		return nil, err
	}

	privKeyBlock, _ := pem.Decode(f)
	if privKeyBlock == nil {
		err = errors.New("invalid private key data")
		return nil, err
	}

	if privKeyBlock.Type == "RSA PRIVATE KEY" {
		rootCertPrivKey, err = x509.ParsePKCS1PrivateKey(privKeyBlock.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		err = errors.New("invalid private key type")
		return nil, err
	}

	rootCertPrivKey.Precompute()

	if err := rootCertPrivKey.Validate(); err != nil {
		return nil, err
	}

	return rootCertPrivKey, nil
}
