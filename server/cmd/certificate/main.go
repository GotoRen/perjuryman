package main

import (
	"fmt"

	"github.com/GotoRen/perjuryman/server/exec"
	"github.com/GotoRen/perjuryman/server/internal"
	"github.com/GotoRen/perjuryman/server/internal/logger"
)

func main() {
	exec.LoadConf()

	// ルートCAを構築
	if err := internal.CreatRootCA(); err != nil {
		logger.LogErr("ルートCAの構築に失敗しました", "error", err)
	} else {
		fmt.Println("[INFO] ルートCAを構築しました")
	}

	// サーバを構築
	if err := internal.CreateServer(); err != nil {
		logger.LogErr("サーバの構築に失敗しました", "error", err)
	} else {
		fmt.Println("[INFO] サーバを構築しました")
	}
}
