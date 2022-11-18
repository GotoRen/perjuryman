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
// サーバを構築します。
func CreateServer() (err error) {
	// ルート証明書とその秘密鍵を取得します。
	rootCert, rootCertPrivKey, err := GetRootCA()
	if err != nil {
		return err
	}

	// サーバのRSA Keyペアを生成
	serverCertPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.LogErr("サーバのRSA Keyペアの生成に失敗しました", "error", err)
		return err
	} else {
		fmt.Println("[INFO] サーバのRSA Keyペアを生成しました")
	}

	// サーバ証明書を発行
	serverCert := setupServerCertificate("server.local") // TODO: env
	if err = generateServerCertificate(serverCert, serverCertPrivKey, rootCert, rootCertPrivKey); err != nil {
		logger.LogErr("サーバ証明書を発行できませんでした。", "error", err)
		return err
	} else {
		fmt.Println("[INFO] サーバ証明書を発行しました")
	}

	// サーバ証明書の秘密鍵を生成
	if err = generateServerCertificatePrivateKey(serverCertPrivKey); err != nil {
		logger.LogErr("サーバ証明書の秘密鍵の生成に失敗しました", "error", err)
		return err
	} else {
		fmt.Println("[INFO] サーバ証明書の秘密鍵を生成しました")
	}

	return nil
}

// サーバのX.509 証明書を表します
func setupServerCertificate(commonName string) (ca *x509.Certificate) {
	var serialNum int64 = 2023 // TODO:  env
	var expandYears int = 10   // TODO:  env

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
		NotAfter:              time.Now().AddDate(expandYears, 0, 0),
		EmailAddresses:        []string{"ren510dev@gmail.com"}, // TODO:  env
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		DNSNames:              []string{commonName},
	}

	return ca
}

// サーバ証明書を発行します
func generateServerCertificate(serverCert *x509.Certificate, serverCertPrivKey *rsa.PrivateKey, rootCert *x509.Certificate, rootCertPrivKey *rsa.PrivateKey) (err error) {
	f, err := os.Create("server.local.pem") // TODO: env
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

// サーバ証明書のRSA秘密鍵を生成します
func generateServerCertificatePrivateKey(serverCertPrivKey *rsa.PrivateKey) (err error) {
	f, err := os.Create("server.local.key") // TODO: env
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
