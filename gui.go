package gui

import (
	"fmt"

	"github.com/rivo/tview"
)

func StartGUI() {
	app := tview.NewApplication()
	sidebar := tview.NewList().
		AddItem("Raum 1", "", '1', nil).
		AddItem("Raum 2", "", '2', nil)

	if err := app.SetRoot(sidebar, true).Run(); err != nil {
		fmt.Println("Fehler GUI:", err)
	}
}
