package p2p

import (
	"fmt"
	"net"
)

func StartServer(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Fehler beim Starten:", err)
		return
	}
	fmt.Println("P2P Server läuft auf", addr)

	for {
		conn, _ := ln.Accept()
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("Nachricht:", string(buf[:n]))
}
