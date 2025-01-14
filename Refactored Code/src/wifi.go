package main

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
	"time"
)

// Function that loads wi-fi sub-menu
func loadWifiMenu() {
	buttons = nil

	sniffButton := tview.NewButton("Sniff").
		SetSelectedFunc(func() {
			pages.SwitchToPage("sniff")
		})

	sniffButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("main") // Switch back to the main page
		})
	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	wifiFlex := tview.NewFlex().
		AddItem(sniffButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	pages.AddPage("wifi", wifiFlex, true, false) // Add the WiFi page to pages
	buttons = append(buttons, sniffButton, backButton)
	enableTabFocus(wifiFlex, buttons)
}

// Function that loads wi-fi sniffing sub-menu
func loadSniffingMenu() {
	buttons = nil

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

	sniffButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").SetSelectedFunc(func() {
		pages.SwitchToPage("wifi")
	})

	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	sniffFlex := tview.NewFlex().
		AddItem(sniffingText, 0, 1, false).
		AddItem(sniffButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = append(buttons, sniffButton, backButton)
	pages.AddPage("sniff", sniffFlex, true, false)
	enableTabFocus(sniffFlex, buttons)
}
