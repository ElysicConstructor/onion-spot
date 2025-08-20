package appgo
package main

import (
	"fmt"
	"github.com/rivo/tview"
)

func startUI(incoming <-chan string, sendFunc func(string) error) {
	app := tview.NewApplication()
	chatView := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() {
		app.Draw()
	})

	input := tview.NewInputField().SetLabel("Nachricht: ")
	input.SetDoneFunc(func(key tview.Key) {
		if key == tview.KeyEnter {
			msg := input.GetText()
			if msg != "" {
				sendFunc(msg)
				fmt.Fprintf(chatView, "[yellow]Ich: %s\n", msg)
				input.SetText("")
			}
		}
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(input, 3, 1, true)

	go func() {
		for msg := range incoming {
			fmt.Fprintf(chatView, "[green]Peer: %s\n", msg)
		}
	}()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
