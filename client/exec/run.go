package exec

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/GotoRen/perjuryman/client/internal"
	"github.com/GotoRen/perjuryman/client/internal/logger"
	"github.com/joho/godotenv"
)

func Run() {
	LoadConf()

	// TLS config
	tlsConf, err := internal.HandleTLS()
	if err != nil {
		logger.LogErr("Failed to get TLS config", "error", err)
	} else {
		fmt.Println("[INFO] Get TLS config")
	}

	// TLS dial
	conn, err := tls.Dial("tcp", os.Getenv("SERVER_FQDN")+":"+os.Getenv("TLS_LISTEN_PORT"), tlsConf)
	if err != nil {
		logger.LogErr("Connection refused", "error", err)
		return
	} else {
		// Connection start
		fmt.Println("[INFO] TLS dial...")
		fmt.Println("[INFO] TLS connection established - Local Addr:", conn.LocalAddr().String())
		fmt.Println("[INFO] TLS connection established - Remote Addr", conn.RemoteAddr().String())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.LogErr("Error when connection closing", "error", err)
		}
	}()

	// Data communication
	internal.PublishMessage(conn)
}

func LoadConf() {
	err := godotenv.Load()
	if err != nil {
		logger.LogErr("Error loading .env file", "error", err)
	}

	logger.InitZap()
}
