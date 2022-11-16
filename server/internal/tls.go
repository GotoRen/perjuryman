package internal

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/GotoRen/perjuryman/server/internal/logger"
)

// GetCert returns serverTLSConf to get server public and private key.
func GetCert() (tlsConf *tls.Config, err error) {
	// サーバ証明書の取得
	// server_cert_name := "./" + os.Getenv("NAME") + ".pem"       // サーバ証明書
	// server_privatekey_name := "./" + os.Getenv("NAME") + ".key" // サーバ証明書の秘密
	server_cert_name := "./server.local.pem"
	server_privatekey_name := "./server.local.key"

	serverCert, err := tls.LoadX509KeyPair(server_cert_name, server_privatekey_name)
	if err != nil {
		logger.LogErr("Failed to load the public and private key pair", "error", err)
		return nil, err
	}

	// ルート認証局情報の取得
	CAPool := x509.NewCertPool()
	if caCert, err := os.ReadFile("ca.pem"); err != nil {
		logger.LogErr("Could not load ca certificate", "error", err)
	} else {
		if ok := CAPool.AppendCertsFromPEM(caCert); !ok {
			err = errors.New("the certificate is not correct")
			logger.LogErr("The certificate is not correct.", "error", err)
			fmt.Println("Error: ルート証明書の解析に失敗しました。") // ca.pemを使用したサーバ証明書を検証
			return nil, err
		}
	}

	// TLS通信のconfig構造を定義
	tlsConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		MinVersion:   tls.VersionTLS13,
		// ClientCAs:    CAPool,
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return tlsConf, nil
}
