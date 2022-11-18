package internal

import (
	"crypto/tls"
	"os"
)

// サーバ証明書を取得します。これはTLSconfigをリターンします。
func HandleTLS() (tlsConf *tls.Config, err error) {
	server_cert_name := os.Getenv("SERVER_CERTIFICATE_NAME")
	server_privatekey_name := os.Getenv("SERVER_CERTIFICATE_PRIVATEKEY_NAME")

	// サーバ証明書を取得
	serverCert, err := tls.LoadX509KeyPair(server_cert_name, server_privatekey_name)
	if err != nil {
		return nil, err
	}

	// TLSconfig構造を定義
	tlsConf = &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{serverCert},
	}

	return tlsConf, nil
}
