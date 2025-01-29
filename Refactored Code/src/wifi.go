package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

	scanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			pages.SwitchToPage("scan")
		})
	scanButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("main") // Switch back to the main page
		})
	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	wifiFlex := tview.NewFlex().
		AddItem(sniffButton, 0, 1, true).
		AddItem(scanButton, 0, 1, false).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	pages.AddPage("wifi", wifiFlex, true, false) // Add the Wi-Fi page to pages
	buttons = []*tview.Button{sniffButton, scanButton, backButton}
	enableTabFocus(wifiFlex, buttons)
}

// Function that loads wi-fi sniffing sub-menu
func loadSniffingMenu() {
	buttons = nil

	sniffingText := tview.NewTextView().
		SetText("Ready to sniff!")
	sniffingText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	sniffButton := tview.NewButton("Start Sniffing").
		SetSelectedFunc(func() {
			logFileName := fmt.Sprintf("sniff_log_%d.pcap", time.Now().Unix())

			cmd := exec.Command("sudo", "ettercap", "-T", "-w", logFileName, "-i", "wlo1")
			stderr := &bytes.Buffer{}
			cmd.Stderr = stderr

			// Start the sniffing process
			if err := cmd.Start(); err != nil {
				sniffingText.SetText(fmt.Sprintf("Error starting command: %v\n", err))
				return
			}

			sniffingText.SetText("Sniffing in progress")

			// Go function that waits for the process to finish and updates the text view
			go func() {
				err := cmd.Wait()
				app.QueueUpdateDraw(func() {
					if err != nil {
						sniffingText.SetText(fmt.Sprintf("Command finished with error: %v\n", err)).
							SetTextColor(tcell.ColorRed)
					} else {
						sniffingText.SetText("Sniffing completed successfully!").
							SetTextColor(tcell.ColorGreen)
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
		AddItem(sniffingText, 0, 3, false).
		AddItem(sniffButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = append(buttons, sniffButton, backButton)
	pages.AddPage("sniff", sniffFlex, true, false)
	enableTabFocus(sniffFlex, buttons)
}

func loadScanMenu() {
	buttons = nil
	var checkESSID, checkAddress, checkProtocol, checkFreq bool

	scanText := tview.NewTextView().
		SetText("Ready to scan network!").
		SetTextColor(tcell.ColorWhite)
	scanText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	ssidCheckbox := tview.NewCheckbox().
		SetLabel("SSID").
		SetChecked(false).
		SetChangedFunc(func(checked bool) {
			checkESSID = checked
		})

	addressCheckbox := tview.NewCheckbox().
		SetLabel("Address").
		SetChecked(false).
		SetChangedFunc(func(checked bool) {
			checkAddress = checked
		})

	protocolCheckbox := tview.NewCheckbox().
		SetLabel("Protocol").
		SetChecked(false).
		SetChangedFunc(func(checked bool) {
			checkProtocol = checked
		})

	frequencyCheckbox := tview.NewCheckbox().
		SetLabel("Frequency").
		SetChecked(false).
		SetChangedFunc(func(checked bool) {
			checkFreq = checked
		})

	checkFlex := tview.NewFlex().
		AddItem(ssidCheckbox, 0, 1, false).
		AddItem(addressCheckbox, 0, 1, false).
		AddItem(protocolCheckbox, 0, 1, false).
		AddItem(frequencyCheckbox, 0, 1, false)

	scanButton := tview.NewButton("Start Scan").
		SetSelectedFunc(func() {
			var essidList, addressList, protocolList, frequencyList, networks, lines []string
			var wg sync.WaitGroup

			scanText.SetText("Scanning for networks...").
				SetTextColor(tcell.ColorWhite)

			cmd := exec.Command("sudo", "iwlist", "wlo1", "scan") // CHANGE THIS TO PI INTERFACE
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				scanText.SetText(fmt.Sprintf("Failed to create pipe: %v\n", err))
				return
			}

			if err := cmd.Start(); err != nil {
				scanText.SetText(fmt.Sprintf("Error starting command: %v\n", err))
				return
			}

			scanner := bufio.NewScanner(stdout)

			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				lines = append(lines, line)
			}

			if checkESSID {
				// Scan line by line for "SSID:" and store the value
				wg.Add(1)
				go func() {
					defer wg.Done()
					for _, line := range lines {
						if strings.HasPrefix(line, "ESSID:") {
							essid := strings.TrimPrefix(line, "ESSID:")
							essidList = append(essidList, strings.TrimSpace(essid))
						}
					}
				}()
			}

			if checkAddress {
				wg.Add(1)
				// Scan line by line for "Address:" and store the value
				go func() {
					defer wg.Done()
					for _, line := range lines {
						if strings.HasPrefix(line, "Cell ") {
							// Split the line at "Address: " to get the MAC address
							parts := strings.Split(line, "Address: ")
							if len(parts) > 1 {
								addressValue := strings.TrimSpace(parts[1])
								// Handle cases where there might be extra text after the MAC address
								addressValue = strings.Split(addressValue, " ")[0]
								addressList = append(addressList, addressValue)
							}
						}
					}
				}()
			}

			if checkProtocol {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for _, line := range lines {
						if strings.HasPrefix(line, "Protocol:") {
							protocolValue := strings.TrimPrefix(line, "Protocol:")
							protocolList = append(protocolList, strings.TrimSpace(protocolValue))
						}
					}
				}()
			}

			if checkFreq {
				wg.Add(1)
				defer wg.Done()
				for _, line := range lines {
					if strings.HasPrefix(line, "Frequency:") {
						frequencyValue := strings.TrimPrefix(line, "Frequency:")
						frequencyList = append(frequencyList, strings.TrimSpace(frequencyValue))
					}
				}
			}

			//Update TUI with list of SSIDs or an error message if none are found
			go func() {
				wg.Wait()
				// Wait for the scan to finish
				err := cmd.Wait()
				if err != nil {

					return
				}

				maxLength := len(essidList)

				if len(addressList) > maxLength {
					maxLength = len(addressList)
				}
				if len(protocolList) > maxLength {
					maxLength = len(protocolList)
				}

				for index := 0; index < maxLength; index++ {
					var thisESSID, thisAddress, thisProtocol, thisFrequency string

					if checkESSID && index < len(essidList) {
						thisESSID = essidList[index]
					}
					if checkAddress && index < len(addressList) {
						thisAddress = addressList[index]
					}
					if checkProtocol && index < len(protocolList) {
						thisProtocol = protocolList[index]
					}

					if checkFreq && index < len(frequencyList) {
						thisFrequency = frequencyList[index]
					}

					thisNetwork := fmt.Sprintf("%s | %s | %s | %s", thisESSID, thisAddress, thisProtocol, thisFrequency)
					networks = append(networks, thisNetwork)
				}

				// Update the scanText with the list of ESSIDs
				app.QueueUpdateDraw(func() {
					if len(networks) > 0 {
						networkList := strings.Join(networks, "\n")
						scanText.SetText(fmt.Sprintf("Found Networks:\n%s", networkList))
					} else {
						scanText.SetText("No networks found.").
							SetTextColor(tcell.ColorRed)
					}
				})
			}()
		})

	scanButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").SetSelectedFunc(func() {
		pages.SwitchToPage("wifi")
	})

	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	scanFlex := tview.NewFlex().
		AddItem(scanText, 0, 3, false).
		AddItem(checkFlex, 0, 1, false).
		AddItem(scanButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = append(buttons, scanButton, backButton)
	pages.AddPage("scan", scanFlex, true, false)
	enableTabFocus(scanFlex, buttons)
}
