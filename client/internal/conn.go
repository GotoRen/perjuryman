package internal

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// HandleTLS establishes TLS connection with the server.
func HandleTLS() (*tls.Conn, error) {
	// ルート認証局の情報を取得
	CAPool := x509.NewCertPool()
	if caCert, err := os.ReadFile("ca.pem"); err != nil {
		return nil, err
	} else {
		if ok := CAPool.AppendCertsFromPEM(caCert); !ok {
			return nil, err
		} else {
			fmt.Println("[INFO] Certificate verified!!")
		}
	}

	// TLS通信のconfig構造を定義
	tlsConf := &tls.Config{
		MinVersion: tls.VersionTLS13,
		RootCAs:    CAPool,
	}

	// TLS dial
	conn, err := tls.Dial("tcp", "server.local:4501", tlsConf)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("[INFO] TLS dial...")
	}

	return conn, nil
}
