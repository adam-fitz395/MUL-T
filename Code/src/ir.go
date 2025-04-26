package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// LoadIRMenu is a function that loads the menu for infrared interactions
func LoadIRMenu() {
	buttons = nil

	scanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			pages.SwitchToPage("irscan")
		})
	scanButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("main")
		})
	backButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	IRFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(scanButton, 0, 1, true).
		AddItem(backButton, 0, 1, false)

	buttons = []*tview.Button{scanButton, backButton}
	pages.AddPage("infrared", IRFlex, true, false)
	EnableTabFocus(IRFlex, buttons)
}

// LoadIRScan is a function that loads the infrared scanning sub-menu
func LoadIRScan() {
	buttons = nil
	scanText := tview.NewTextView().
		SetDynamicColors(true).SetText("[green]Ready to scan!")
	scanText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	scanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			go func() {
				app.QueueUpdateDraw(func() {
					scanText.SetText("[yellow]Stopping LIRC services...")
				})

				// Stop LIRC service
				stopCmd := exec.Command("sudo", "systemctl", "stop", "lircd")
				if err := stopCmd.Run(); err != nil {
					// Check if error is acceptable "not active" state
					var exitErr *exec.ExitError
					if errors.As(err, &exitErr) {
						if exitErr.ExitCode() == 5 { // systemctl exit code 5 = unit not loaded
							// This is acceptable, continue
						} else {
							app.QueueUpdateDraw(func() {
								scanText.SetText("[red]Error stopping service: " + err.Error())
							})
							return
						}
					}
				}

				// Change infrared interface
				changeInterface := []*exec.Cmd{
					// Set driver to default
					exec.Command("sudo", "sed", "-i",
						"s/^driver.*/driver = default/",
						"/etc/lirc/lirc_options.conf"),

					// Set device to /dev/lirc1 (note escaped slashes)
					exec.Command("sudo", "sed", "-i",
						"s/^device.*/device = \\/dev\\/lirc1/",
						"/etc/lirc/lirc_options.conf"),
				}

				for _, cmd := range changeInterface {
					if err := cmd.Run(); err != nil {
						app.QueueUpdateDraw(func() {
							scanText.SetText("[red]Error changing interface: " + err.Error())
						})
						return
					}
				}

				// Capture IR signal
				app.QueueUpdateDraw(func() {
					scanText.SetText("[green]Press remote button now... (10 second window)")
				})

				cmd := exec.Command("timeout", "5", "mode2", "-d", "/dev/lirc1")
				output, err := cmd.CombinedOutput()
				if err != nil {
					// Ignress timeout error (expected)
					if !strings.Contains(err.Error(), "exit status 124") {
						app.QueueUpdateDraw(func() {
							scanText.SetText("[red]Capture error: " + err.Error())
						})
						return
					}
				}

				// Process timings
				app.QueueUpdateDraw(func() {
					scanText.SetText("[yellow]Processing timings...")
				})

				var timings []string
				scanner := bufio.NewScanner(bytes.NewReader(output))
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "pulse ") || strings.HasPrefix(line, "space ") {
						parts := strings.Fields(line)
						if len(parts) == 2 {
							timings = append(timings, parts[1])
						}
					}
				}

				// Validate results
				if len(timings) < 20 {
					app.QueueUpdateDraw(func() {
						scanText.SetText(fmt.Sprintf(
							"[red]Insufficient data (%d timings)", len(timings)))
					})
					return
				}

				// Fix odd count
				if len(timings)%2 != 0 {
					timings = timings[:len(timings)-1]
				}

				// Create directory if needed
				dirPath := "../logfiles/rawir"
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Error creating directory: " + err.Error())
					})
					return
				}

				// Generate filename with timestamp
				filename := filepath.Join(dirPath, fmt.Sprintf("ir_%s.txt",
					time.Now().Format("20060102_150405"))) // YYYYMMDD_HHMMSS format

				// Save to file
				content := strings.Join(timings, " ")
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Save failed: " + err.Error())
					})
					return
				}

				app.QueueUpdateDraw(func() {
					scanText.SetText(fmt.Sprintf(
						"[green]Captured %d timings\nSaved to [white]ir_timings.txt",
						len(timings)))
				})
			}()
		})
	scanButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	scanBackButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("infrared")
		})
	scanBackButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	scanFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(scanText, 0, 3, false).
		AddItem(scanButton, 0, 1, true).
		AddItem(scanBackButton, 0, 1, false)

	buttons = []*tview.Button{scanButton, scanBackButton}
	pages.AddPage("irscan", scanFlex, true, false)
	EnableTabFocus(scanFlex, buttons)
}
