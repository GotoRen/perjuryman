package main

import (
	"fmt"

	"github.com/GotoRen/perjuryman/server/exec"
	"github.com/GotoRen/perjuryman/server/internal"
	"github.com/GotoRen/perjuryman/server/internal/logger"
)

func main() {
	exec.LoadConf()

	// built the RootCA
	if err := internal.CreatRootCA(); err != nil {
		logger.LogErr("Failed to build RootCA", "error", err)
	} else {
		fmt.Println("[INFO] Built the RootCA")
	}

	// built the Server
	if err := internal.CreateServer(); err != nil {
		logger.LogErr("Failed to build Server", "error", err)
	} else {
		fmt.Println("[INFO] Built the Server")
	}
}
