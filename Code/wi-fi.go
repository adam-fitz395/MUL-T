package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os/exec"
	"time"
)

var logFileName string
var isSniffing bool

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
		outputLabel, // Label for sniffing output

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

	isSniffing = true
	outputLabel.Hidden = false

	// Create a temporary log file with a Unix timestamp
	logFileName = fmt.Sprintf("sniff_log_%d.pcap", time.Now().Unix())

	// Define command to run ettercap on devices network interface and to log using the defined logfile
	cmd = exec.Command("sudo", "ettercap", "-T", "-w", logFileName, "-i", "wlo1") // This line will need to change to be the wireless interface on the Pi!!!

	// Create a buffered stderr for error capture
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	// Start the sniffing process
	if err := cmd.Start(); err != nil {
		log.Println("Error starting command:", err)
		updateLabel(outputLabel, "Error starting sniffing: "+err.Error())
		return
	}

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

	isSniffing = false
}
