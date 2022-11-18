package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/GotoRen/perjuryman/server/internal/logger"
)

func SubscribeMessage(conn net.Conn) {
	for {
		// 受信
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.LogErr("The connection has been broken", "error", err)
			return
		}
		fmt.Print("\n[DEBUG] Message Received: " + string(message))

		// 送信
		newmessage := strings.ToUpper(message)
		_, err = conn.Write([]byte(newmessage + "\n"))
		if err != nil {
			logger.LogErr("Failed to write packet", "error", err)
			return
		}
		fmt.Print("[DEBUG] Will send reply message: " + string(newmessage))
	}
}
