package internal

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// サーバ証明書を取得します。これはTLSconfigをリターンします。
func GetServverCertificate() (tlsConf *tls.Config, err error) {
	server_cert_name := os.Getenv("SERVER_CERTIFICATE_NAME")
	server_privatekey_name := os.Getenv("SERVER_CERTIFICATE_PRIVATEKEY_NAME")

	// サーバ証明書を取得
	serverCert, err := tls.LoadX509KeyPair(server_cert_name, server_privatekey_name)
	if err != nil {
		return nil, err
	}

	// ルート認証局情報の取得
	caPool := x509.NewCertPool()
	if pemCert, err := os.ReadFile(os.Getenv("ROOT_CERTIFICATE_NAME")); err != nil {
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
