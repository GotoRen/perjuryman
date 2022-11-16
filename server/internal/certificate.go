package internal

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
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
// Create server certificate
// ==================================================================//

// GenerateServerCert creates server certificate.
func GenerateServerCert() (tlsConf *tls.Config, err error) {
	// set up RootCA
	ca, caPrivKey, err := GetRootCA()
	if err != nil {
		return nil, err
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	cert := SetUpServerCA(os.Getenv("NAME"))

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	certPEM, err := CreateCert(certBytes)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("[INFO] Create server certificate")
	}

	tlsConf, err = CreateCertPrivate(certPrivKey, certPEM)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("[INFO] Create server Private Key")
	}

	return tlsConf, nil
}

// SetUpServerCA is setting up ServerCA.
func SetUpServerCA(commonName string) (ca *x509.Certificate) {
	var serialNum int64 = 2023
	var expandYears int = 11

	// set up our ServerCA certificate
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
		BasicConstraintsValid: true,
		IsCA:                  true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(expandYears, 0, 0), // 10 years
		EmailAddresses:        []string{os.Getenv("SERVER_EMAIL")},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		DNSNames:              []string{commonName},
	}
	return ca
}

// CreateCert creates server ServerCA.
func CreateCert(certBytes []byte) (certPEM *bytes.Buffer, err error) {
	f, err := os.Create(os.Getenv("NAME") + ".pem")
	if err != nil {
		logger.LogErr("server: Failed to create server.pem", "error", err)
		return nil, err
	}

	// pem encode
	certPEM = new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		logger.LogErr("server: Failed to create Buffer string", "error", err)
		return nil, err
	}

	err = pem.Encode(f, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		logger.LogErr("server: server.pem encoding failure", "error", err)
		return nil, err
	}

	err = f.Close()
	if err != nil {
		logger.LogErr("server: Failed to close server.pem", "error", err)
		return nil, err
	}

	return certPEM, nil
}

// CreateCertPrivate creates server private certificates.
func CreateCertPrivate(certPrivKey *rsa.PrivateKey, certPEM *bytes.Buffer) (tlsConf *tls.Config, err error) {
	// サーバの秘密鍵を作成
	f, err := os.Create(os.Getenv("NAME") + ".key")
	if err != nil {
		logger.LogErr("server-private: Failed to create server.pem", "error", err)
		return nil, err
	}

	// pem encode
	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		logger.LogErr("server-private: Failed to create Buffer string", "error", err)
		return nil, err
	}

	err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		logger.LogErr("server-private: server.pem encoding failure", "error", err)
		return nil, err
	}

	err = f.Close()
	if err != nil {
		logger.LogErr("server-private: Failed to close server.pem", "error", err)
		return nil, err
	}

	// TLS通信のセットアップ
	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		logger.LogErr("server-private: Failed to parse the public and private key pair.", "error", err)
		return nil, err
	}

	tlsConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	return tlsConf, nil
}

// // ==================================================================//
// // Create client certificate
// // ==================================================================//

// // CreateClientCert creates client certificates.
// func CreateClientCert(certBytes []byte, cName string) (err error) {
// 	f, err := os.Create(cName + "_cert.pem")
// 	if err != nil {
// 		logger.LogErr("client: Failed to create cName_cert.pem", "error", err)
// 		return err
// 	}

// 	err = pem.Encode(f, &pem.Block{
// 		Type:  "CERTIFICATE",
// 		Bytes: certBytes,
// 	})
// 	if err != nil {
// 		logger.LogErr("client: cName_cert.pem encoding failure", "error", err)
// 		return err
// 	}

// 	err = f.Close()
// 	if err != nil {
// 		logger.LogErr("client: Failed to close cName_cert.pem", "error", err)
// 		return err
// 	}

// 	return nil
// }
