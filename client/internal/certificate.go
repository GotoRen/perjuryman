// Package certificate contains the certficate of TLS connetction.
package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"

	"github.com/GotoRen/perjuryman/client/internal/logger"
)

// GenerateCertRequest return cyphonic client certificate request.
func GenerateCertRequest(cName string, nEmail string) (ccertBytes []byte, err error) {
	privateKey, err := CreateRSAKey()
	if err != nil {
		logger.LogErr("Failed to get the key information", "error", err)
	}
	client := SetUpClientCertRequest(cName, nEmail, &privateKey.PublicKey)

	csr, err := x509.CreateCertificateRequest(rand.Reader, client, privateKey)
	if err != nil {
		return nil, err
	}

	err = CreateClientCertRequest(csr, cName)
	if err != nil {
		return nil, err
	}

	return csr, nil
}

// SetUpClientCertRequest is setting up client certificate request struct.
func SetUpClientCertRequest(commonName string, nodeEmail string, publicKey *rsa.PublicKey) (client *x509.CertificateRequest) {
	// set up client certificate request
	client = &x509.CertificateRequest{
		PublicKeyAlgorithm: x509.RSA,
		PublicKey:          publicKey,
		SignatureAlgorithm: x509.SHA256WithRSA,
		Subject: pkix.Name{
			Organization:       []string{"Perjuryman"},
			OrganizationalUnit: []string{"Perjuryman"},
			Country:            []string{"JP"},
			Province:           []string{"Tokyo"},
			Locality:           []string{"Shibuya"},
			StreetAddress:      []string{""},
			PostalCode:         []string{""},
			CommonName:         commonName,
		},
		EmailAddresses: []string{nodeEmail},
		DNSNames:       []string{commonName},
	}
	return client
}

// CreateClientCertRequest creates client certificate request.
func CreateClientCertRequest(certBytes []byte, cName string) (err error) {
	f, err := os.Create(cName + "_certreq.pem")
	if err != nil {
		logger.LogErr("client: Failed to create cName_certreq.pem", "error", err)
		return err
	}

	err = pem.Encode(f, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: certBytes,
	})
	if err != nil {
		logger.LogErr("client: cName_certreq.pem encoding failure", "error", err)
		return err
	}

	err = f.Close()
	if err != nil {
		logger.LogErr("client: Failed to close cName_certreq.pem", "error", err)
		return err
	}

	return nil
}
