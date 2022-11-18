package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/GotoRen/perjuryman/server/internal/logger"
)

// ==================================================================//
// Create "Server certificate"
// ==================================================================//

// CreateServer builds a Server.
func CreateServer() (err error) {
	// get the root certificate and its private key
	rootCert, rootCertPrivKey, err := GetRootCA()
	if err != nil {
		return err
	}

	// generate Server's rsa key pair
	serverCertPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.LogErr("Failed to generate the Server's rsa key pair", "error", err)
		return err
	} else {
		fmt.Println("[INFO] Generate the Server's rsa key pair")
	}

	// issue the server certificate
	serverCert := setupServerCertificate(os.Getenv("SERVER_FQDN")) // Define the Server's X.509 certificate
	if err = generateServerCertificate(serverCert, serverCertPrivKey, rootCert, rootCertPrivKey); err != nil {
		logger.LogErr("Failed to generate the server certificate", "error", err)
		return err
	} else {
		fmt.Println("[INFO] Generate the server certificate")
	}

	// generate private key for server certificate
	if err = generateServerCertificatePrivateKey(serverCertPrivKey); err != nil {
		logger.LogErr("Failed to generate private key for server certificate", "error", err)
		return err
	} else {
		fmt.Println("[INFO] Generate private key for server certificate")
	}

	return nil
}

// setupServerCertificate represents the Server's X.509 certificate.
func setupServerCertificate(commonName string) (cert *x509.Certificate) {
	var serialNum int64 = 2023
	var expandYears int = 10

	cert = &x509.Certificate{
		SerialNumber: big.NewInt(serialNum),
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
		BasicConstraintsValid: true,
		IsCA:                  true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(expandYears, 0, 0),
		EmailAddresses:        []string{os.Getenv("CERTIFICATE_REGISTER_EMAIL")},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		DNSNames:              []string{commonName},
	}

	return cert
}

// generateServerCertificate generates server certificate.
func generateServerCertificate(serverCert *x509.Certificate, serverCertPrivKey *rsa.PrivateKey, rootCert *x509.Certificate, rootCertPrivKey *rsa.PrivateKey) (err error) {
	f, err := os.Create(os.Getenv("SERVER_CERTIFICATE_NAME"))
	if err != nil {
		return err
	}

	b, err := x509.CreateCertificate(rand.Reader, serverCert, rootCert, &serverCertPrivKey.PublicKey, rootCertPrivKey)
	if err != nil {
		return err
	}

	// pem encode
	if err = pem.Encode(f, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: b,
	}); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil
}

// generateServerCertificatePrivateKey generates RSA private key for server certificate.
func generateServerCertificatePrivateKey(serverCertPrivKey *rsa.PrivateKey) (err error) {
	f, err := os.Create(os.Getenv("SERVER_CERTIFICATE_PRIVATEKEY_NAME"))
	if err != nil {
		return err
	}

	// pem encode
	if err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverCertPrivKey),
	}); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil
}
