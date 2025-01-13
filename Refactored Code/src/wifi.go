package main

import (
	"bytes"
	"fmt"
	"github.com/rivo/tview"
	"os/exec"
	"time"
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
			logFileName := fmt.Sprintf("sniff_log_%d.pcap", time.Now().Unix())

			cmd := exec.Command("sudo", "ettercap", "-T", "-w", logFileName, "-i", "eth0")
			stderr := &bytes.Buffer{}
			cmd.Stderr = stderr

			// Start the sniffing process
			if err := cmd.Start(); err != nil {
				sniffingText.SetText(fmt.Sprintf("Error starting command: %v\n", err))
				return
			}

			// Go function that waits for the process to finish and updates the text view
			go func() {
				err := cmd.Wait()
				app.QueueUpdateDraw(func() {
					if err != nil {
						sniffingText.SetText(fmt.Sprintf("Command finished with error: %v\n", err))
					} else {
						sniffingText.SetText("Sniffing completed successfully!")
					}
				})
			}()
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
