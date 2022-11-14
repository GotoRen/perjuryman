package exec

import (
	"fmt"

	"github.com/GotoRen/perjuryman/client/internal/logger"
	"github.com/joho/godotenv"
)

func Run() {
	LoadConf()

	fmt.Println("Hello, World!")
}

func LoadConf() {
	err := godotenv.Load()
	if err != nil {
		logger.LogErr("Error loading .env file", "error", err)
	}

	logger.InitZap()
}
