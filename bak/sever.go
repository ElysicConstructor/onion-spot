package appgo
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

var key = []byte("thisis32bytepasswordforaesgcm!!!") // 32 Byte AES Key

func startServer(port string, incoming chan<- string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Fehler beim Starten des Servers:", err)
		return
	}
	fmt.Println("Warte auf Verbindungen auf Port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		remote := conn.RemoteAddr().String()
		fmt.Printf("Verbindungsanfrage von %s akzeptieren? (y/n): ", remote)

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		if answer != "y" {
			fmt.Println("Abgelehnt.")
			conn.Close()
			continue
		}

		// Push-Notification
		exec.Command("notify-send", "Secure P2P Chat", "Neue Verbindung von "+remote).Run()

		go handleConnection(conn, incoming)
	}
}

func handleConnection(conn net.Conn, incoming chan<- string) {
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		msg, err := decryptMessage(key, buf[:n])
		if err == nil {
			incoming <- string(msg)
		}
	}
}

func sendMessage(addr, msg string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	enc, err := encryptMessage(key, []byte(msg))
	if err != nil {
		return err
	}
	_, err = conn.Write(enc)
	return err
}
