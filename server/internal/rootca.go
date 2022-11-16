package internal

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/GotoRen/perjuryman/server/internal/logger"
)

// ==================================================================//
// Create CA / Root certificate
// ==================================================================//

// GenerateRootCA creates RootCA.
func GenerateRootCA() (ca *x509.Certificate, caPrivKey *rsa.PrivateKey, err error) {
	ca = SetUpRootCA("RootCA")

	// create our private and public key
	caPrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	if err = CreateCA(ca, caPrivKey); err != nil {
		return nil, nil, err
	} else {
		fmt.Println("[INFO] Create RootCA's certificate")
	}

	if err = CreateCAPrivateKeyPEM(ca, caPrivKey); err != nil {
		return nil, nil, err
	} else {
		fmt.Println("[INFO] Create RootCA's Private Key")
	}

	return ca, caPrivKey, nil
}

// SetUpRootCA is setting up RootCA.
func SetUpRootCA(commonName string) (ca *x509.Certificate) {
	var serialNum int64 = 2023
	var expandYears int = 11

	// set up our RootCA certificate
	ca = &x509.Certificate{
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
		NotAfter:              time.Now().AddDate(expandYears, 0, 0), // 10 years
		IsCA:                  true,
		EmailAddresses:        []string{os.Getenv("ROOTCA_EMAIL")},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	return ca
}

// CreateCA creates RootCA.
func CreateCA(ca *x509.Certificate, caPrivKey *rsa.PrivateKey) (err error) {
	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		logger.LogErr("RootCA: x.509 certificate creation failed", "error", err)
		return err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		logger.LogErr("RootCA: Failed to create Buffer string", "error", err)
		return err
	}

	// create root certificate
	f, err := os.Create("ca.pem")
	if err != nil {
		logger.LogErr("RootCA: Failed to create ca.pem", "error", err)
		return err
	}

	err = pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes})
	if err != nil {
		logger.LogErr("RootCA: ca.pem encoding failure", "error", err)
		return err
	}

	err = f.Close()
	if err != nil {
		logger.LogErr("RootCA: Failed to close ca.pem", "error", err)
		return err
	}

	return nil
}

// CreateCAPrivateKeyPEM creates Private Key.
func CreateCAPrivateKeyPEM(ca *x509.Certificate, caPrivKey *rsa.PrivateKey) (err error) {
	// ルート認証局の秘密鍵を作成
	f, err := os.Create("ca_private_key.pem")
	if err != nil {
		logger.LogErr("RootCA-private: Failed to create ca_private_key.pem", "error", err)
		return err
	}

	// pem encode
	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		logger.LogErr("RootCA-private: Failed to create Buffer string", "error", err)
		return err
	}

	err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		logger.LogErr("RootCA-private: ca_private_key.pem encoding failure", "error", err)
		return err
	}

	err = f.Close()
	if err != nil {
		logger.LogErr("RootCA-private: Failed to close ca_private_key.pem", "error", err)
		return err
	}

	return nil
}

// ==================================================================//
// Read CA / Root certificate
// ==================================================================//

// GetRootCA read ca cert and caPrivKey
func GetRootCA() (ca *x509.Certificate, caPrivKey *rsa.PrivateKey, err error) {
	caPrivKey, err = ReadCAPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	ca, err = ReadCACertificate()
	if err != nil {
		return nil, nil, err
	}

	return ca, caPrivKey, nil
}

// ReadCAPrivateKey read ca PrivKey
func ReadCAPrivateKey() (caPrivKey *rsa.PrivateKey, err error) {
	PriKey, err := os.ReadFile("ca_private_key.pem")
	if err != nil {
		logger.LogErr("RootCA: Could not load ca privtekey", "error", err)
		return nil, err
	}

	PrivKeyBlock, _ := pem.Decode(PriKey)
	if PrivKeyBlock == nil {
		err = errors.New("invalid private key data")
		logger.LogErr("RootCA: private key data faild", "error", err)
		return nil, err
	}

	if PrivKeyBlock.Type == "RSA PRIVATE KEY" {
		caPrivKey, err = x509.ParsePKCS1PrivateKey(PrivKeyBlock.Bytes)
		if err != nil {
			logger.LogErr("RootCA: Failed to parse private key", "error", err)
			return nil, err
		}
	} else {
		err = errors.New("invalid private key type")
		logger.LogErr("RootCA: Error private key type", "error", err)
		return nil, err
	}

	caPrivKey.Precompute()

	if err := caPrivKey.Validate(); err != nil {
		logger.LogErr("RootCA: Error Integrity of private key", "error", err)
		return nil, err
	}

	return caPrivKey, nil
}

// ReadCACertificate read ca Certificate
func ReadCACertificate() (ca *x509.Certificate, err error) {
	Cert, err := os.ReadFile("ca.pem")
	if err != nil {
		logger.LogErr("RootCA: Could not load ca certificate", "error", err)
		return nil, err
	}
	CAPool := x509.NewCertPool()
	if ok := CAPool.AppendCertsFromPEM(Cert); !ok {
		logger.LogErr("RootCA: The certificate is not correct.")
		return nil, err
	}

	CertBlock, _ := pem.Decode(Cert)
	if CertBlock == nil {
		err = errors.New("invalid private key data")
		logger.LogErr("RootCA: ca certificate data faild", "error", err)
		return nil, err
	}
	ca, err = x509.ParseCertificate(CertBlock.Bytes)
	if err != nil {
		logger.LogErr("RootCA: failed to parse certificate", "error", err)
		return nil, err
	}

	return ca, nil
}
