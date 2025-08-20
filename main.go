package main

import (
	"fmt"
	"mychatapp/gui"
	"mychatapp/p2p"
)

func main() {
	fmt.Println("Starte MyChatApp...")
	go p2p.StartServer(":6666") // P2P Server starten
	gui.StartGUI()              // GUI starten
}
