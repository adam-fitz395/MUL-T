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
			loadBluetoothScanningMenu(multiWindow, outputLabel)
		}),
		widget.NewButton("Show Devices", func() {
		}),
		widget.NewButton("Back", func() {
			loadMainMenu(multiWindow, outputLabel)
		}),
	))
}

func loadBluetoothScanningMenu(multiWindow fyne.Window, outputLabel *widget.Label) {
	multiWindow.SetContent(container.NewVBox(
		outputLabel,
		widget.NewButton("Scan", func() {
			startBluetoothScanning(outputLabel)
		}),
		widget.NewButton("Back", func() {
			stopBluetoothScanning(outputLabel)
			loadBluetoothMenu(multiWindow, outputLabel)
		}),
	))
}

func startBluetoothScanning(outputLabel *widget.Label) {
	outputLabel.Hidden = false
	updateLabel(outputLabel, "This is only a test, nothing's happening currently ¯\\_(O_o)_/¯")
}
func stopBluetoothScanning(outputLabel *widget.Label) {

}
