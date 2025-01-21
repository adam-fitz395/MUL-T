package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
	"strings"
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

	pages.AddPage("wifi", wifiFlex, true, false) // Add the WiFi page to pages
	buttons = []*tview.Button{sniffButton, scanButton, backButton}
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

			cmd := exec.Command("sudo", "ettercap", "-T", "-w", logFileName, "-i", "wlo1")
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

func loadScanMenu() {
	buttons = nil

	var checkSSID, checkLastSeen bool

	scanText := tview.NewTextView().
		SetText("Ready to scan network!").
		SetTextColor(tcell.ColorWhite)
	scanText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	ssidCheckbox := tview.NewCheckbox().
		SetLabel("SSID").
		SetChecked(false).
		SetChangedFunc(func(checked bool) {
			checkSSID = checked
		})

	lastSeenCheckbox := tview.NewCheckbox().
		SetLabel("Last seen").
		SetChecked(false).
		SetChangedFunc(func(checked bool) {
			checkLastSeen = checked
		})

	checkFlex := tview.NewFlex().
		AddItem(ssidCheckbox, 0, 1, false).
		AddItem(lastSeenCheckbox, 0, 1, false)

	scannerForm := tview.NewForm().
		AddCheckbox("SSID", false, func(checked bool) {
			checkSSID = checked
		}).
		AddCheckbox("Last seen", false, func(checked bool) {
			checkLastSeen = checked
		})

	scannerForm.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	scanButton := tview.NewButton("Start Scan").
		SetSelectedFunc(func() {
			var ssids, lastSeenList, networks []string
			scanText.SetText("Scanning for networks...").
				SetTextColor(tcell.ColorWhite)

			cmd := exec.Command("sudo", "iw", "wlo1", "scan") // CHANGE THIS TO PI INTERFACE
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				scanText.SetText(fmt.Sprintf("Failed to create pipe: %v\n", err))
				return
			}

			if err := cmd.Start(); err != nil {
				scanText.SetText(fmt.Sprintf("Error starting command: %v\n", err))
				return
			}

			// Channel to handle concurrent scanning and output to TUI
			done := make(chan struct{})

			if checkSSID == true {
				// Scan line by line for "SSID:" and store the value
				go func() {
					scanner := bufio.NewScanner(stdout)
					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())
						if strings.HasPrefix(line, "SSID:") {
							ssid := strings.TrimPrefix(line, "SSID:")
							ssids = append(ssids, strings.TrimSpace(ssid))
						}
					}

					// Handle potential scanner error
					if err := scanner.Err(); err != nil {
						app.QueueUpdateDraw(func() {
							scanText.SetText(fmt.Sprintf("Error reading output: %v\n", err))
						})
						return
					}

					// Signal scan completion
					done <- struct{}{}
				}()
			}

			if checkLastSeen == true {
				// Scan line by line for "last seen:" and store the value
				go func() {
					scanner := bufio.NewScanner(stdout)
					for scanner.Scan() {
						line := strings.TrimSpace(scanner.Text())
						if strings.HasPrefix(line, "last seen:") {
							lastSeenValue := strings.TrimPrefix(line, "last seen:")
							lastSeenList = append(lastSeenList, strings.TrimSpace(lastSeenValue))
						}
					}

					// Handle potential scanner error
					if err := scanner.Err(); err != nil {
						app.QueueUpdateDraw(func() {
							scanText.SetText(fmt.Sprintf("Error reading output: %v\n", err))
						})
						return
					}

					// Signal scan completion
					done <- struct{}{}
				}()
			}

			//Update TUI with list of SSIDs or an error message if none are found
			go func() {
				<-done
				// Wait for the scan to finish
				cmd.Wait()

				for index := range ssids {
					var thisNetwork string
					var thisSSID, thisLastSeen string

					if checkSSID && index < len(ssids) {
						thisSSID = ssids[index]
					}
					if checkLastSeen && index < len(lastSeenList) {
						thisLastSeen = lastSeenList[index]
					}

					thisNetwork = fmt.Sprintf("%s | %s", thisSSID, thisLastSeen)
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
		AddItem(scanText, 0, 1, false).
		AddItem(checkFlex, 0, 1, false).
		AddItem(scanButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = append(buttons, scanButton, backButton)
	pages.AddPage("scan", scanFlex, true, false)
	enableTabFocus(scanFlex, buttons)
}
