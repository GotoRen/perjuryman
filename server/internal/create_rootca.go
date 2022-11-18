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
// ルートCAを構築します。
func CreatRootCA() (err error) {
	// ルートCAのRSA Keyペアを生成
	rootCertPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.LogErr("ルートCAのRSA Keyペアの生成に失敗しました", "error", err)
		return err
	} else {
		fmt.Println("[INFO] ルートCAのRSA Keyペアを生成しました")
	}

	// ルート証明書を発行
	rootCert := setupRootCertificate("RootCA") // ルートCAのX.509 証明書を定義します。
	if err = generateRootCertificate(rootCert, rootCertPrivKey); err != nil {
		logger.LogErr("ルート証明書を発行できませんでした。", "error", err)
		return err
	} else {
		fmt.Println("[INFO] ルート証明書を発行しました")
	}

	// ルート証明書の秘密鍵を生成
	if err = generateRootCertificatePrivateKey(rootCertPrivKey); err != nil {
		logger.LogErr("ルート証明書の秘密鍵の生成に失敗しました", "error", err)
		return err
	} else {
		fmt.Println("[INFO] ルート証明書の秘密鍵を生成しました")
	}

	return nil
}

// ルートCAのX.509 証明書を表します
func setupRootCertificate(commonName string) (ca *x509.Certificate) {
	var serialNum int64 = 2023
	var expandYears int = 10

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
		NotAfter:              time.Now().AddDate(expandYears, 0, 0),
		IsCA:                  true,
		EmailAddresses:        []string{os.Getenv("CERTIFICATE_REGISTER_EMAIL")},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	return ca
}

// ルート証明書を発行します
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

// ルート証明書のRSA秘密鍵を生成します
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
