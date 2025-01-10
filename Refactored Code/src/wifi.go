package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/rivo/tview"
)

// Function that loads wi-fi sub-menu
func loadWifiMenu() {
	sniffButton := tview.NewButton("Sniff").
		SetSelectedFunc(func() {
			pages.SwitchToPage("sniff")
		})

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("main") // Switch back to the main page
		})

	wifiFlex := tview.NewFlex().
		AddItem(sniffButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	pages.AddPage("wifi", wifiFlex, true, false) // Add the WiFi page to pages
}

// Function that loads wi-fi sniffing sub-menu
func loadSniffingMenu() {
	sniffingText := tview.NewTextView().
		SetText("Ready to sniff!")

	sniffButton := tview.NewButton("Start Sniffing").
		SetSelectedFunc(func() {
			logFileName := fmt.Sprintf("sniff_log_%d.pcap", time.Now().Unix()) // Set logfile name using current time

			cmd := exec.Command("sudo", "ettercap", "-T", "-w", logFileName, "-i", "eth0") // This line will need to change to be the wireless interface on the Pi!!!

			// Create a buffered stderr for error capture
			stderr := &bytes.Buffer{}
			cmd.Stderr = stderr

			sniffingText.Clear()
			sniffingText.SetText("Sniffing in Progress!")

			// Start the sniffing process
			if err := cmd.Start(); err != nil {
				log.Println("Error starting command:", err)
				return
			}

			if err := cmd.Wait(); err != nil {
				log.Println("Command finished with error:", err)
			}
		})

	backButton := tview.NewButton("Back").SetSelectedFunc(func() {
		pages.SwitchToPage("wifi")
	})

	sniffFlex := tview.NewFlex().
		AddItem(sniffingText, 0, 1, false).
		AddItem(sniffButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	pages.AddPage("sniff", sniffFlex, true, false)
}
