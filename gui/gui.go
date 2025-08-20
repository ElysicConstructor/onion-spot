package gui

import (
	"fmt"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func StartTUI(name, room string, conn *net.UDPConn) error {
	peers := NewPeerSet()
	app := tview.NewApplication()

	sidebar := tview.NewList().ShowSecondaryText(false)
	sidebar.SetBorder(true).SetTitle("Räume")
	sidebar.AddItem(room, "", 0, nil)

	chatBox := tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
	chatBox.SetBorder(true).SetTitle("Chat")

	input := tview.NewInputField().SetLabel("Nachricht: ").SetFieldWidth(0)
	input.SetBorder(true)

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			msg := input.GetText()
			if msg != "" {
				line := fmt.Sprintf("%s: %s", name, msg)
				chatBox.Write([]byte(fmt.Sprintf("[yellow]%s\n", line)))
				input.SetText("")
				for _, p := range peers.List() {
					conn.WriteToUDP([]byte("MSG "+line), p)
				}
			}
		}
	})

	grid := tview.NewGrid().
		SetRows(0, 3).
		SetColumns(30, 0).
		SetBorders(true)
	grid.AddItem(sidebar, 0, 0, 2, 1, 0, 0, true)
	grid.AddItem(chatBox, 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(input, 1, 1, 1, 1, 0, 0, true)

	app.SetRoot(grid, true).SetFocus(input)

	// Nachricht-Empfang
	go func() {
		buf := make([]byte, 4096)
		for {
			n, from, _ := conn.ReadFromUDP(buf)
			line := strings.TrimSpace(string(buf[:n]))
			if strings.HasPrefix(line, "MSG ") {
				line = strings.TrimPrefix(line, "MSG ")
				chatBox.Write([]byte(fmt.Sprintf("[green]%s\n", line)))
				app.Draw()
			} else if strings.HasPrefix(line, "PUNCH") {
				peers.Add(from)
			}
		}
	}()

	return app.Run()
}

// ---------------- PeerSet für GUI ----------------
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
	return out
}
