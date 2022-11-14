package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// CreateRSAKey creates client rsakey.
func CreateRSAKey() (key *rsa.PrivateKey, err error) {
	nName := os.Getenv("NAME")

	rsaPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	err = createPrivate(rsaPrivKey, nName)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("[INFO] Create client private key.")
	}

	err = createPublic(&rsaPrivKey.PublicKey, nName)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("[INFO] Create client public key.")
	}

	return rsaPrivKey, nil
}

// createCertPrivate creates private.key.
func createPrivate(certPrivKey *rsa.PrivateKey, nName string) (err error) {
	f, err := os.Create(nName + "_private.key")
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// createCertPrivate creates private key.
func createPublic(certPublicKey *rsa.PublicKey, nName string) (err error) {
	f, err := os.Create(nName + ".pem")
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = pem.Encode(f, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(certPublicKey),
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
