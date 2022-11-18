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

	// TLS config
	tlsConf, err := internal.HandleTLS()
	if err != nil {
		logger.LogErr("Failed to get TLS config", "error", err)
	} else {
		fmt.Println("[INFO] Get TLS config")
	}

	// Listen
	ln, err := tls.Listen("tcp", ":"+os.Getenv("TLS_LISTEN_PORT"), tlsConf)
	if err != nil {
		logger.LogErr("Connection refused", "error", err)
		return
	} else {
		fmt.Println("[INFO] TLS listen...")
	}
	defer func() {
		if err := ln.Close(); err != nil {
			logger.LogErr("Error when TLS listen closing", "error", err)
		}
	}()

	// Connection start
	conn, err := ln.Accept()
	if err != nil {
		logger.LogErr("Can't get the socket", "error", err)
	} else {
		fmt.Println("[INFO] TLS connection established - Local Addr:", conn.LocalAddr().String())
		fmt.Println("[INFO] TLS connection established - Remote Addr", conn.RemoteAddr().String())
	}
	defer conn.Close()

	// Data communication
	internal.SubscribeMessage(conn)
}

func LoadConf() {
	err := godotenv.Load()
	if err != nil {
		logger.LogErr("Error loading .env file", "error", err)
	}

	logger.InitZap()
}
