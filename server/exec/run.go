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

	serverTLSConf, err := internal.GetCert() // サーバ証明書を取得
	if err != nil {
		logger.LogErr("Failed to issue perjuryman server certificate", "error", err)
	}

	// TLS listen
	ln, err := tls.Listen("tcp", "server.local:"+os.Getenv("TLS_PORT"), serverTLSConf)
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

	conn, err := ln.Accept()
	if err != nil {
		logger.LogErr("Can't get the socket", "error", err)
	}
	defer conn.Close()

	fmt.Println("[INFO] TLS connection established - Local Addr:", conn.LocalAddr().String())
	fmt.Println("[INFO] TLS connection established - Remote Addr", conn.RemoteAddr().String())

	internal.SubscribeMessage(conn)
}

func LoadConf() {
	err := godotenv.Load()
	if err != nil {
		logger.LogErr("Error loading .env file", "error", err)
	}

	logger.InitZap()
}
