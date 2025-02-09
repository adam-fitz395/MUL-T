package main

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
	"strings"
	"sync"
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
		SetLabel("MAC Address").
		SetChecked(false).
		SetChangedFunc(func(checked bool) {
			checkMAC = checked
		})

	btCheckFlex := tview.NewFlex().
		AddItem(MACCheckbox, 0, 1, false)

	btScanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			var lines, macList, devices []string
			wg := new(sync.WaitGroup)

			btScanText.SetText("Scanning for Bluetooth Devices...")

			cmd := exec.Command("sudo", "timeout", "5", "bluetoothctl", "scan", "on")
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

			if checkMAC {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for _, line := range lines {
						if strings.HasPrefix(line, "[NEW] Device ") {
							mac := strings.TrimPrefix(line, "[NEW] Device ")
							macList = append(macList, strings.TrimSpace(mac))
						}
					}
				}()
			}

			go func() {
				wg.Wait()
				// Wait for the scan to finish
				err := cmd.Wait()
				if err != nil {

					return
				}

				maxLength := len(macList)

				for index := 0; index < maxLength; index++ {
					var thisMAC string

					if checkMAC && index < len(macList) {
						thisMAC = macList[index]
					}

					thisDevice := fmt.Sprintf("%s", thisMAC)
					devices = append(devices, thisDevice)
				}

				// Update the scanText with the list of MACs
				app.QueueUpdateDraw(func() {
					if len(devices) > 0 {
						deviceList := strings.Join(devices, "\n")
						btScanText.SetText(fmt.Sprintf("Found Devices:\n%s", deviceList))
					} else {
						btScanText.SetText("[red]No networks found.[red]")
					}
				})
			}()

			//TODO: Display bluetooth device information
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
