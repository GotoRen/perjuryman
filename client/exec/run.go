package exec

import (
	"github.com/GotoRen/perjuryman/client/internal"
	"github.com/GotoRen/perjuryman/client/internal/logger"
	"github.com/joho/godotenv"
)

func Run() {
	LoadConf()

	if _, err := internal.HandleTLS(); err != nil {
		logger.LogErr("Failed to establish TLS connection", "error", err)
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
