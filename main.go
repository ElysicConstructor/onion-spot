package main

import (
	"fmt"
	"time"

	"github.com/ElysicConstructor/onion-spot/p2p"
)

func main() {
	// Name / Raum
	name := fmt.Sprintf("peer-%d", time.Now().Unix()%10000)
	room := "default"

	// Automatischer Start: überprüft, ob ein Introducer im Netzwerk läuft
	p2p.AutoStart(name, room)
}
