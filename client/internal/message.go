package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/GotoRen/perjuryman/client/internal/logger"
)

func PublishMessage(conn net.Conn) {
Loop:
	for {
		// Read in input from stdin
		text := getInput()

		// If the input is "exit", we will stop the loop and end client.
		switch {
		case text == "exit\n":
			log.Print("Entering exit command...")
			break Loop
		}

		// Send
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			logger.LogErr("Failed to write packet", "error", err)
			return
		}
		fmt.Print("[DEBUG] Will send text: " + text)

		// Receive
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.LogErr("Unable to retrieve message", "error", err)
			return
		}
		fmt.Print("[DEBUG] Message Received: " + string(message))
	}

	fmt.Println("[DEBUG] Successfully closed the TLS connection.")
}

func getInput() string {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nInput text: ")
		text, _ := reader.ReadString('\n')

		if len(text) > 1 {
			return text
		}
	}
}
