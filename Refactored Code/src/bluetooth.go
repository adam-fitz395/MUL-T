package main

import (
	"bufio"
	"context"
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
	var checkMAC bool

	buttons = nil
	btScanText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]Ready to Scan!")

	btScanText.SetBorder(true)

	MACCheckbox := tview.NewCheckbox().
		SetLabel("MAC Address ").
		SetChecked(false).
		SetChangedFunc(func(checked bool) {
			checkMAC = checked
		})

	btCheckFlex := tview.NewFlex().
		AddItem(MACCheckbox, 0, 1, false)

	btScanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			btScanText.SetText("[white] Scanning in progress...please wait!")
			go func() {
				var macList, devices []string

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				cmd := exec.CommandContext(ctx, "sudo", "bluetoothctl", "--", "scan", "on")
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					app.QueueUpdateDraw(func() {
						btScanText.SetText(fmt.Sprintf("[red]Failed to create pipe: %v\n[red]", err))
					})
					return
				}

				if err := cmd.Start(); err != nil {
					app.QueueUpdateDraw(func() {
						btScanText.SetText(fmt.Sprintf("[red]Error starting command: %v\n[red]", err))
					})
					return
				}

				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())

					// Check for MAC addresses in real-time
					if checkMAC {
						if strings.HasPrefix(line, "[NEW] Device ") {
							macValue := strings.TrimPrefix(line, "[NEW] Device ")
							macList = append(macList, macValue)
						}
					}
				}

				// Process collected MAC addresses
				for _, mac := range macList {
					devices = append(devices, mac)
				}

				// Update the UI
				app.QueueUpdateDraw(func() {
					if len(devices) > 0 {
						deviceList := strings.Join(devices, "\n")
						btScanText.SetText(fmt.Sprintf("Found Devices:\n%s", deviceList))
					} else {
						btScanText.SetText("[red]No devices found.[red]")
					}
				})
			}()
		})

	backButton := tview.NewButton("Back").SetSelectedFunc(func() {
		pages.SwitchToPage("bluetooth")
	})

	btScanFlex := tview.NewFlex().
		AddItem(btScanText, 0, 1, false).
		AddItem(btCheckFlex, 0, 1, false).
		AddItem(btScanButton, 0, 1, true).
		AddItem(backButton, 0, 1, false)
	btScanFlex.SetDirection(tview.FlexRow)

	buttons = []*tview.Button{btScanButton, backButton}
	pages.AddPage("btScan", btScanFlex, true, false) // Add the Wi-Fi page to pages
	enableTabFocus(btScanFlex, buttons)
}
