package internal

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// サーバ証明書を取得します。これはTLSconfigをリターンします。
func GetServverCertificate() (tlsConf *tls.Config, err error) {
	server_cert_name := "server.local.pem"
	server_privatekey_name := "server.local.key"

	// サーバ証明書を取得
	serverCert, err := tls.LoadX509KeyPair(server_cert_name, server_privatekey_name)
	if err != nil {
		return nil, err
	}

	// ルート認証局情報の取得
	caPool := x509.NewCertPool()
	if pemCert, err := os.ReadFile("ca.pem"); err != nil {
		return nil, err
	} else {
		if ok := caPool.AppendCertsFromPEM(pemCert); !ok {
			err = errors.New("The certificate is incorrect")
			return nil, err
		}
	}

	// TLSconfig構造を定義
	tlsConf = &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{serverCert},
	}

	return tlsConf, nil
}
