package exec

import (
	"fmt"

	"github.com/GotoRen/perjuryman/client/internal"
	"github.com/GotoRen/perjuryman/client/internal/logger"
	"github.com/joho/godotenv"
)

func Run() {
	LoadConf()

	conn, err := internal.HandleTLS()
	if err != nil {
		logger.LogErr("Failed to establish TLS connection...", "error", err)
	} else {
		fmt.Println("TLSを確立しました:", conn.LocalAddr().String())
		fmt.Println("TLSを確立しました:", conn.RemoteAddr().String())
	}

	for {
	}

	// _ := HandleServerConn(conn)
}

func LoadConf() {
	err := godotenv.Load()
	if err != nil {
		logger.LogErr("Error loading .env file", "error", err)
	}

	logger.InitZap()
}
