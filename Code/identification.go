package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func loadNFCMenu(multiWindow fyne.Window, outputLabel *widget.Label) {
	multiWindow.SetContent(container.NewVBox(
		widget.NewButton("Scan", func() {
		}),
		widget.NewButton("Emulate", func() {
		}),
		widget.NewButton("Back", func() {
			loadMainMenu(multiWindow, outputLabel)
		}),
	))
}
