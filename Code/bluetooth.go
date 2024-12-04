package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var isScanning bool

func loadBluetoothMenu(multiWindow fyne.Window, outputLabel *widget.Label) {
	multiWindow.SetContent(container.NewVBox(
		widget.NewButton("Scan for Bluetooth devices", func() {
		}),
		widget.NewButton("Back", func() {
			loadMainMenu(multiWindow, outputLabel)
		}),
	))
}

func startBluetoothScanning(outputLabel *widget.Label) {

}
func stopBluetoothScanning(outputLabel *widget.Label) {

}
