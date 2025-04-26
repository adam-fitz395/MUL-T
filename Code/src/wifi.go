package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
	"strings"
	"sync"
)

// LoadWifiMenu is a function that loads the Wi-Fi sub-menu
func LoadWifiMenu() {
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

	MITMButton := tview.NewButton("MITM").
		SetSelectedFunc(func() {
			pages.SwitchToPage("mitm")
		})
	MITMButton.SetBorder(true).
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
		AddItem(MITMButton, 0, 1, false).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	pages.AddPage("wifi", wifiFlex, true, false) // Add the Wi-Fi page to pages
	buttons = []*tview.Button{sniffButton, scanButton, MITMButton, backButton}
	EnableTabFocus(wifiFlex, buttons)
}

// LoadSniffingMenu is a function that loads the wi-fi sniffing sub-menu
func LoadSniffingMenu() {
	buttons = nil
	var duration int

	sniffingText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]Ready to sniff![green]")

	sniffingText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	sniffDuration := tview.NewDropDown().SetLabel("Duration (Seconds): ").SetLabelColor(tcell.ColorWhite)
	sniffDuration.
		AddOption("10",
			func() {
				duration = 10
			}).
		AddOption("20", func() {
			duration = 20
		}).
		AddOption("30", func() {
			duration = 30
		}).AddOption("60", func() {
		duration = 60
	})
	sniffDuration.SetCurrentOption(0)

	sniffButton := tview.NewButton("Start Sniffing").
		SetSelectedFunc(func() {
			sniffingText.SetText("[white]Sniffing in progress...please wait!")
			go func() {
				cmd := exec.Command("bash", "../scripts/wifi_sniff.sh", fmt.Sprintf("%d", duration))

				// Start the script to initiate scanning
				if err := cmd.Start(); err != nil {
					app.QueueUpdateDraw(func() {
						sniffingText.SetText(fmt.Sprintf("[red]Error starting script:[white] %v\n", err))
					})
					return
				}

				// Go function that waits for the process to finish and updates the text view
				err := cmd.Wait()
				if err != nil {
					app.QueueUpdateDraw(func() {
						sniffingText.SetText(fmt.Sprintf("[red]Script execution error:[white] %v\n", err))
					})
					return
				}
				app.QueueUpdateDraw(func() {
					if err != nil {
						sniffingText.SetText(fmt.Sprintf("[red]Command finished with error:[white] %v\n", err))
					} else {
						sniffingText.SetText("[green]Sniffing completed successfully!\n[blue]A pcap file has been created!")
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
		AddItem(sniffDuration, 0, 1, false).
		AddItem(sniffButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = append(buttons, sniffButton, backButton)
	pages.AddPage("sniff", sniffFlex, true, false)
	EnableTabFocus(sniffFlex, buttons)
}

// LoadScanMenu is a function that loads the wi-fi scan sub-menu
func LoadScanMenu() {
	buttons = nil
	var checkESSID, checkAddress, checkFreq = true, true, true

	scanText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]Ready to scan network![green]").
		SetTextColor(tcell.ColorWhite)

	scanText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	ssidCheckbox := tview.NewCheckbox().
		SetLabel("SSID").
		SetChecked(true).
		SetChangedFunc(func(checked bool) {
			checkESSID = checked
		})

	addressCheckbox := tview.NewCheckbox().
		SetLabel("Address").
		SetChecked(true).
		SetChangedFunc(func(checked bool) {
			checkAddress = checked
		})

	frequencyCheckbox := tview.NewCheckbox().
		SetLabel("Frequency").
		SetChecked(true).
		SetChangedFunc(func(checked bool) {
			checkFreq = checked
		})

	checkFlex := tview.NewFlex().
		AddItem(ssidCheckbox, 0, 1, false).
		AddItem(addressCheckbox, 0, 1, false).
		AddItem(frequencyCheckbox, 0, 1, false)

	scanButton := tview.NewButton("Start Scan").
		SetSelectedFunc(func() {
			var essidList, addressList, frequencyList, networks, lines []string
			var wg sync.WaitGroup

			scanText.SetText("[white]Scanning for networks...")

			go func() {
				cmd := exec.Command("sudo", "iwlist", "wlan0", "scan") // CHANGE THIS TO PI INTERFACE
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					scanText.SetText(fmt.Sprintf("[red]Failed to create pipe: %v\n[red]", err))
					return
				}

				if err := cmd.Start(); err != nil {
					scanText.SetText(fmt.Sprintf("[red]Error starting command: %v\n[red]", err))
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

					for index := 0; index < maxLength; index++ {
						var thisESSID, thisAddress, thisFrequency string

						if checkESSID && index < len(essidList) {
							thisESSID = essidList[index]
						}

						if checkAddress && index < len(addressList) {
							thisAddress = addressList[index]
						}

						if checkFreq && index < len(frequencyList) {
							thisFrequency = frequencyList[index]
						}

						thisNetwork := fmt.Sprintf("%s | %s | %s", thisESSID, thisAddress, thisFrequency)
						networks = append(networks, thisNetwork)
					}

					// Update the scanText with the list of ESSIDs
					app.QueueUpdateDraw(func() {
						if len(networks) > 0 {
							networkList := strings.Join(networks, "\n")
							scanText.SetText(fmt.Sprintf("Found Networks:\n%s", networkList))
						} else {
							scanText.SetText("[red]No networks found.[red]")
						}
					})
				}()
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

	buttons = []*tview.Button{scanButton, backButton}
	pages.AddPage("scan", scanFlex, true, false)
	EnableTabFocus(scanFlex, buttons)
}

// LoadMITMMenu is a function that loads the man-in-the-middle attack sub-menu
func LoadMITMMenu() {
	MITMText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]Ready to perform MITM attack!")

	MITMText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	startMITMButton := tview.NewButton("Start").
		SetSelectedFunc(func() {
			MITMText.SetText("[red]MITM attack in progress! Press [blue]Stop[red] to stop attack!")
			go func() {
				cmd := exec.Command("bash", "../scripts/mitm.sh", "On")

				// Capture command output
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &out

				if err := cmd.Start(); err != nil {
					app.QueueUpdateDraw(func() {
						MITMText.SetText(fmt.Sprintf("[red]Error starting MITM: %v", err))
					})
					return
				}

				// Wait for command completion in goroutine
				go func() {
					err := cmd.Wait()
					app.QueueUpdateDraw(func() {
						if err != nil {
							MITMText.SetText(fmt.Sprintf(
								"[red]MITM failed:[white]\n%s\n[red]Error: %v",
								out.String(), err,
							))
						} else {
							MITMText.SetText(fmt.Sprintf(
								"[green]MITM started successfully!\n[white]%s",
								out.String(),
							))
						}
					})
				}()
			}()
		})

	startMITMButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	stopMITMButton := tview.NewButton("Stop").SetSelectedFunc(func() {
		cmd := exec.Command("bash", "../scripts/mitm.sh", "Off")

		// Capture command output
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Start(); err != nil {
			app.QueueUpdateDraw(func() {
				MITMText.SetText(fmt.Sprintf("[red]Error stopping MITM: %v", err))
			})
			return
		}

		// Handle command completion
		go func() {
			err := cmd.Wait()
			app.QueueUpdateDraw(func() {
				if err != nil {
					MITMText.SetText(fmt.Sprintf(
						"[red]Stop failed:[white]\n%s\n[red]Error: %v",
						out.String(), err,
					))
				} else {
					MITMText.SetText(fmt.Sprintf(
						"[green]MITM stopped!\n[blue]Log created!\n[white]%s",
						out.String(),
					))
				}
			})
		}()
	})

	stopMITMButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	MITMBackButton := tview.NewButton("Back").SetSelectedFunc(func() {
		pages.SwitchToPage("wifi")
	})
	MITMBackButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	MITMFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(MITMText, 0, 3, false).
		AddItem(startMITMButton, 0, 1, true).
		AddItem(stopMITMButton, 0, 1, false).
		AddItem(MITMBackButton, 0, 1, false)

	buttons = []*tview.Button{startMITMButton, stopMITMButton, MITMBackButton}
	pages.AddPage("mitm", MITMFlex, true, false)
	EnableTabFocus(MITMFlex, buttons)
}
