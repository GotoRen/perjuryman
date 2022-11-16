package internal

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

// HandleTLS establishes TLS connection with the server.
func HandleTLS() (*tls.Conn, error) {
	// ルート認証局の情報を取得
	CAPool := x509.NewCertPool()
	if caCert, err := os.ReadFile("ca.pem"); err != nil {
		return nil, err
	} else {
		if ok := CAPool.AppendCertsFromPEM(caCert); !ok {
			return nil, err
		} else {
			fmt.Println("[INFO] Certificate verified!!")
		}
	}

	// TLS通信のconfig構造を定義
	tlsConf := &tls.Config{
		RootCAs: CAPool,
	}

	// TLS通信要求
	conn, err := tls.Dial("tcp", "server.local:4501", tlsConf)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func PublishMessage(conn net.Conn) {
Loop:
	for {
		// Read in input from stdin
		text := GetInput()

		// If the input is "exit", we will stop the loop and end client.
		switch {
		case text == "exit\n":
			log.Print("Entering exit command. EXIT.")
			break Loop
		}
		log.Print("Will send text: " + text)

		// Send to socket
		conn.Write([]byte(text + "\n"))

		// Listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Message from server: " + message)
	}
}

func GetInput() string {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Input text: ")
		text, _ := reader.ReadString('\n')

		if len(text) > 1 {
			return text
		}
	}
}
