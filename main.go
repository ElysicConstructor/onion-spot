package appgo
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  Server: go run . server <port>")
		fmt.Println("  Client: go run . client <ip:port>")
		return
	}

	mode := os.Args[1]
	addr := os.Args[2]
	incoming := make(chan string)

	if mode == "server" {
		go startServer(addr, incoming)
		startUI(incoming, func(msg string) error {
			return nil // Server sendet nicht direkt
		})
	} else if mode == "client" {
		startUI(incoming, func(msg string) error {
			return sendMessage(addr, msg)
		})
	}
}
