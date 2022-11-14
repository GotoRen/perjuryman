package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GotoRen/perjuryman/client/exec"
	"github.com/GotoRen/perjuryman/client/internal"
)

func main() {
	exec.LoadConf()

	_, err := internal.GenerateCertRequest(os.Getenv("NAME"), os.Getenv("EMAIL"))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("[INFO] Create client certificate.")
	}
}
