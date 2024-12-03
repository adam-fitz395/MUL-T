package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os/exec"
	"time"
)

var cmd *exec.Cmd
var logFileName string
var isSniffing bool

func main() {
	// Create app and window
	multiMenu := app.New()
	multiWindow := multiMenu.NewWindow("Network Sniffer")

	// Create a label with word wrapping
	outputLabel := widget.NewLabel("")
	outputLabel.Wrapping = fyne.TextWrapWord

	loadMainMenu(multiWindow, outputLabel)

	// Resize window and show
	multiWindow.Resize(fyne.NewSize(400, 300)) // Set a reasonable window size
	multiWindow.ShowAndRun()
}

// Function to load main menu
func loadMainMenu(multiWindow fyne.Window, outputLabel *widget.Label) {
	multiWindow.SetContent(container.NewVBox(
		widget.NewButton("Wi-Fi", func() {
			loadWifiMenu(multiWindow, outputLabel)
		}),

		widget.NewButton("Bluetooth", func() {
			// Insert sub-menu function
		}),

		widget.NewButton("NFC", func() {
			// Insert sub-menu function
		}),
	))
}

// Function that loads Wi-Fi sub-menu
func loadWifiMenu(multiWindow fyne.Window, outputLabel *widget.Label) {
	multiWindow.SetContent(container.NewVBox(
		widget.NewButton("Ettercap Sniff", func() {
			loadWifiSniffer(multiWindow, outputLabel)
		}),
		widget.NewButton("Back", func() {
			loadMainMenu(multiWindow, outputLabel)
		}),
	))
}

// Function to load wi-fi sniffer sub-menu
func loadWifiSniffer(multiWindow fyne.Window, outputLabel *widget.Label) {
	multiWindow.SetContent(container.NewVBox(
		outputLabel, // Show the label for output
		widget.NewButton("Start Sniffing", func() {
			outputLabel.SetText("Starting sniffing...")
			go runSniffer(outputLabel)
		}),
		widget.NewButton("Stop Sniffing", func() {
			go stopSniffing(outputLabel)
		}),
		widget.NewButton("Back", func() {
			if isSniffing == true {
				stopSniffing(outputLabel) // Stop sniffing if it is still ongoing
			}
			loadWifiMenu(multiWindow, outputLabel)
		}),
	))
}

// Function that runs wi-fi sniffer when start button is hit
func runSniffer(outputLabel *widget.Label) {

	// Set sniffing tracker to true
	isSniffing = true

	// Create a temporary log file with a timestamp
	logFileName = fmt.Sprintf("sniff_log_%d.pcap", time.Now().Unix())

	// Set command to run ettercap on devices interface and to log using the defined logfile
	cmd = exec.Command("sudo", "ettercap", "-T", "-w", logFileName, "-i", "wlo1")

	// Create a buffered stderr for error capture
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	// Start the sniffing process
	if err := cmd.Start(); err != nil {
		log.Println("Error starting command:", err)
		updateLabel(outputLabel, "Error starting sniffing: "+err.Error())
		return
	}

	// Display that sniffing is in progress
	updateLabel(outputLabel, "Sniffing in progress...")

	// Wait for the command to finish (without printing stdout)
	if err := cmd.Wait(); err != nil {
		log.Println("Command finished with error:", err)
		updateLabel(outputLabel, "Command finished with error: "+err.Error()+"\n"+stderr.String())
	}
}

// Function that stops sniffing when stop button is hit
func stopSniffing(outputLabel *widget.Label) {
	// Check if there is an ongoing sniffing process to stop
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			log.Println("Error stopping command:", err)
			updateLabel(outputLabel, "Error stopping sniffing.")
			return
		}
		// After stopping, display the filename of the saved log
		updateLabel(outputLabel, fmt.Sprintf("Sniffing stopped. Log saved as %s", logFileName))
	} else {
		updateLabel(outputLabel, "No sniffing process to stop.")
	}

	// Set sniffing tracker to false
	isSniffing = false
}

// Update label to display different text
func updateLabel(label *widget.Label, s string) {
	label.SetText(s)
	label.Refresh()
}
