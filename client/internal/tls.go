package internal

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// TLSconfigを返します。
func HandleTLS() (tlsConf *tls.Config, err error) {
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
		MinVersion: tls.VersionTLS13,
		RootCAs:    caPool,
	}

	return tlsConf, nil
}
