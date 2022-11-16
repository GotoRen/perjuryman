package exec

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/GotoRen/perjuryman/server/internal"
	"github.com/GotoRen/perjuryman/server/internal/logger"
	"github.com/joho/godotenv"
)

func Run() {
	LoadConf()

	fmt.Println("Hello, World!")

	serverTLSConf, err := internal.GetCert()
	if err != nil {
		logger.LogErr("Failed to issue perjuryman server certificate", "error", err)
	}
	fmt.Println("OK")

	ln, err := tls.Listen("tcp", ":"+os.Getenv("TLS_PORT"), serverTLSConf)
	if err != nil {
		logger.LogErr("Connection refused", "error", err)
		return
	}

	defer func() {
		if err := ln.Close(); err != nil {
			logger.LogErr("Error when TLS listen closing", "error", err)
		}
	}()

	for {
		_, err := ln.Accept()
		if err != nil {
			logger.LogErr("Can't get the socket", "error", err)
			continue
		} else {
			fmt.Println("[INFO] Established TLS connection...")
		}

		// go internal.HandleConnection(conn, db)
	}
}

func LoadConf() {
	err := godotenv.Load()
	if err != nil {
		logger.LogErr("Error loading .env file", "error", err)
	}

	logger.InitZap()
}
