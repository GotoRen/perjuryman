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
		logger.LogErr("Failed to establish TLS connection", "error", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.LogErr("Error when connection closing", "error", err)
		}
	}()

	fmt.Println("[INFO] TLS connection established - Local Addr:", conn.LocalAddr().String())
	fmt.Println("[INFO] TLS connection established - Remote Addr", conn.RemoteAddr().String())

	// HandleConnection(conn)
	internal.PublishMessage(conn)
	// RoutineSequentialSender()
	// RoutineSequentialReceiver()
}

func LoadConf() {
	err := godotenv.Load()
	if err != nil {
		logger.LogErr("Error loading .env file", "error", err)
	}

	logger.InitZap()
}
