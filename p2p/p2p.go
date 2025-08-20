package p2p

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/ElysicConstructor/onion-spot/gui"
)

const DefaultPort = 5555
const BroadcastAddr = "255.255.255.255"

// ---------------- PeerSet ----------------
type PeerSet struct {
	addrs map[string]*net.UDPAddr
}

func NewPeerSet() *PeerSet {
	return &PeerSet{addrs: make(map[string]*net.UDPAddr)}
}

func (ps *PeerSet) Add(addr *net.UDPAddr) {
	ps.addrs[addr.String()] = addr
}

func (ps *PeerSet) List() []*net.UDPAddr {
	out := make([]*net.UDPAddr, 0, len(ps.addrs))
	for _, a := range ps.addrs {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].String() < out[j].String() })
	return out
}

// ---------------- Introducer ----------------
type RoomState struct {
	peers map[string]*net.UDPAddr
}

type Introducer struct {
	rooms map[string]*RoomState
	conn  *net.UDPConn
}

func RunIntroducer(listen string) error {
	addr, _ := net.ResolveUDPAddr("udp", listen)
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	intr := &Introducer{rooms: make(map[string]*RoomState), conn: conn}
	fmt.Println("Introducer läuft auf", conn.LocalAddr())

	buf := make([]byte, 2048)
	for {
		n, from, _ := conn.ReadFromUDP(buf)
		msg := strings.TrimSpace(string(buf[:n]))
		go intr.handle(from, msg)
	}
}

func (in *Introducer) handle(from *net.UDPAddr, msg string) {
	if msg == "HELLO" {
		in.reply(from, "I_AM_INTRODUCER")
		return
	}

	parts := strings.Fields(msg)
	if len(parts) < 2 {
		return
	}
	cmd := strings.ToUpper(parts[0])
	switch cmd {
	case "JOIN":
		if len(parts) < 3 {
			return
		}
		room := parts[1]
		name := parts[2]
		rs := in.rooms[room]
		if rs == nil {
			rs = &RoomState{peers: make(map[string]*net.UDPAddr)}
			in.rooms[room] = rs
		}
		rs.peers[from.String()] = from
		fmt.Println("Peer beigetreten:", name, from.String())
	case "LEAVE":
		room := parts[1]
		if rs, ok := in.rooms[room]; ok {
			delete(rs.peers, from.String())
			if len(rs.peers) == 0 {
				delete(in.rooms, room)
			}
		}
	}
}

func (in *Introducer) reply(to *net.UDPAddr, msg string) {
	_, _ = in.conn.WriteToUDP([]byte(msg), to)
}

// ---------------- Peer TUI ----------------
func RunPeerTUI(name, room, introducerAddr, listen string) error {
	laddr, _ := net.ResolveUDPAddr("udp", listen)
	conn, _ := net.ListenUDP("udp", laddr)
	defer conn.Close()

	return gui.StartTUI(name, room, conn)
}

// ---------------- Auto-Start ----------------
func DetectIntroducer() *net.UDPAddr {
	laddr, _ := net.ResolveUDPAddr("udp", ":0")
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Println("Fehler beim Öffnen des Sockets:", err)
		return nil
	}
	defer conn.Close()

	raddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", BroadcastAddr, DefaultPort))
	_, _ = conn.WriteToUDP([]byte("HELLO"), raddr)

	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buf := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil
	}
	if string(buf[:n]) == "I_AM_INTRODUCER" {
		return addr
	}
	return nil
}

func AutoStart(name, room string) {
	introducerAddr := DetectIntroducer()
	if introducerAddr != nil {
		fmt.Println("Introducer gefunden:", introducerAddr.String())
		RunPeerTUI(name, room, introducerAddr.String(), fmt.Sprintf(":%d", DefaultPort))
	} else {
		fmt.Println("Kein Introducer gefunden. Starte selbst als Introducer.")
		RunIntroducer(fmt.Sprintf(":%d", DefaultPort))
	}
}
