package exec

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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

	// TLS通信ダイアル
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

	for {
		// _, err := ln.Accept()
		// if err != nil {
		// 	logger.LogErr("Can't get the socket", "error", err)
		// 	continue
		// } else {
		// 	fmt.Println("[INFO] Established TLS connection...")
		// }

		// go HandleConnection(conn)
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("TLS接続成功")
		}

		_, err = conn.Write([]byte("hello\n"))
		if err != nil {
			fmt.Println("[-] ERROR:", err)
			return
		} else {
			fmt.Println("Send messages")
		}

		// defer conn.Close()

		// for {
		// 	message, err := bufio.NewReader(conn).ReadString('\n')
		// 	if err != nil {
		// 		fmt.Println("no read", err)
		// 	}
		// 	log.Print("Message Received:", string(message))
		// 	newmessage := strings.ToUpper(message)
		// 	conn.Write([]byte(newmessage + "\n"))
		// }
	}
}

func LoadConf() {
	err := godotenv.Load()
	if err != nil {
		logger.LogErr("Error loading .env file", "error", err)
	}

	logger.InitZap()
}

// HandleConnection handle TCP connection.
func HandleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			logger.LogErr("Error when connection closing", "error", err)
		}
	}()

	for {
		b := make([]byte, 3500)
		bi := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		_, err := bi.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			return
		}
		fmt.Println("OK:", b)
	}
}
