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
// Create "RootCA" and "Root certificate"
// ==================================================================//

// CreateRootCA builds a RootCA.
func CreatRootCA() (err error) {
	// generate RootCA's rsa key pair
	rootCertPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.LogErr("Failed to generate the RootCA's rsa key pair", "error", err)
		return err
	} else {
		fmt.Println("[INFO] Generate the RootCA's rsa key pair")
	}

	// issue the root certificate
	rootCert := setupRootCertificate("RootCA") // Define the RootCA's X.509 certificate
	if err = generateRootCertificate(rootCert, rootCertPrivKey); err != nil {
		logger.LogErr("Failed to generate the root certificate", "error", err)
		return err
	} else {
		fmt.Println("[INFO] Generate the root certificate")
	}

	// generate private key for root certificate
	if err = generateRootCertificatePrivateKey(rootCertPrivKey); err != nil {
		logger.LogErr("Failed to generate private key for root certificate", "error", err)
		return err
	} else {
		fmt.Println("[INFO] Generate private key for root certificate")
	}

	return nil
}

// setupRootCertificate represents the RootCA's X.509 certificate.
func setupRootCertificate(commonName string) (cert *x509.Certificate) {
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
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(expandYears, 0, 0),
		IsCA:                  true,
		EmailAddresses:        []string{os.Getenv("CERTIFICATE_REGISTER_EMAIL")},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	return cert
}

// generateRootCertificate generates root certificate.
func generateRootCertificate(rootCert *x509.Certificate, rootCertPrivKey *rsa.PrivateKey) (err error) {
	f, err := os.Create(os.Getenv("ROOT_CERTIFICATE_NAME"))
	if err != nil {
		return err
	}

	b, err := x509.CreateCertificate(rand.Reader, rootCert, rootCert, &rootCertPrivKey.PublicKey, rootCertPrivKey)
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

// generateRootCertificatePrivateKey generates RSA private key for root certificate.
func generateRootCertificatePrivateKey(rootCertPrivKey *rsa.PrivateKey) (err error) {
	f, err := os.Create(os.Getenv("ROOT_CERTIFICATE_PRIVATEKEY_NAME"))
	if err != nil {
		return err
	}

	// pem encode
	if err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rootCertPrivKey),
	}); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil
}
