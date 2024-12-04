package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"os/exec"
)

var cmd *exec.Cmd

func main() {
	// Create app and window
	multiMenu := app.New()
	multiWindow := multiMenu.NewWindow("Network Sniffer")
	multiWindow.Resize(fyne.NewSize(400, 300)) // Set a reasonable window size

	// Create a label with word wrapping and hide
	outputLabel := widget.NewLabel("")
	outputLabel.Wrapping = fyne.TextWrapWord
	outputLabel.Hidden = true

	loadMainMenu(multiWindow, outputLabel)
	multiWindow.ShowAndRun()
}

// Function to load main menu
func loadMainMenu(multiWindow fyne.Window, outputLabel *widget.Label) {
	multiWindow.SetContent(container.NewVBox(
		widget.NewButton("Wi-Fi", func() {
			loadWifiMenu(multiWindow, outputLabel)
		}),

		widget.NewButton("Bluetooth", func() {
			loadBluetoothMenu(multiWindow, outputLabel)
		}),

		widget.NewButton("NFC/RFID", func() {
			loadNFCMenu(multiWindow, outputLabel)
		}),

		widget.NewButton("IR", func() {
			// Insert sub-menu function
		}),

		// The button below will either be removed or functionally changed when implemented onto Pi,
		// for now it just quits the application
		widget.NewButton("Quit", func() {
			os.Exit(0)
		}),
	))
}

// Update label text with new value
func updateLabel(label *widget.Label, s string) {
	label.SetText(s)
	label.Refresh()
}
