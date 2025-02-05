package main

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
	"strings"
	"time"
)

func loadBluetoothMenu() {
	buttons = nil

	btScanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			pages.SwitchToPage("btScan")
		})
	btScanButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("main") // Switch back to the main page
		})
	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	bluetoothFlex := tview.NewFlex().
		AddItem(btScanButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = []*tview.Button{btScanButton, backButton}
	pages.AddPage("bluetooth", bluetoothFlex, true, false) // Add the Wi-Fi page to pages
	enableTabFocus(bluetoothFlex, buttons)
}

func loadBluetoothScan() {
	buttons = nil
	btScanText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]Ready to Scan!")

	btScanText.SetBorder(true)

	btScanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			var lines []string
			btScanText.SetText("Scanning for Bluetooth Devices...")

			cmd := exec.Command("sudo", "bluetoothctl", "--", "scan", "on")
			time.Sleep(5 * time.Second)
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				btScanText.SetText(fmt.Sprintf("[red]Failed to create pipe: %v\n[red]", err))
				return
			}

			if err := cmd.Start(); err != nil {
				btScanText.SetText(fmt.Sprintf("[red]Error starting command: %v\n[red]", err))
				return
			}

			scanner := bufio.NewScanner(stdout)

			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				lines = append(lines, line)
			}

			//TODO: Display bluetooth device information
		})

	backButton := tview.NewButton("Back").SetSelectedFunc(func() {
		pages.SwitchToPage("bluetooth")
	})

	btScanFlex := tview.NewFlex().
		AddItem(btScanText, 0, 1, false).
		AddItem(btScanButton, 0, 1, true).
		AddItem(backButton, 0, 1, false)
	btScanFlex.SetDirection(tview.FlexRow)

	buttons = []*tview.Button{btScanButton, backButton}
	pages.AddPage("btScan", btScanFlex, true, false) // Add the Wi-Fi page to pages
	enableTabFocus(btScanFlex, buttons)
}
