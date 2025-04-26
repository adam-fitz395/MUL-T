package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"os/exec"
	"strings"
)

// Function that loads the menu for infrared functions
func loadIRMenu() {
	buttons = nil

	scanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			pages.SwitchToPage("irscan")
		})
	scanButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	emitButton := tview.NewButton("Emit").
		SetSelectedFunc(func() {
			pages.SwitchToPage("iremit")
		})
	emitButton.SetBorder(true).
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
		AddItem(backButton, 0, 1, false).
		AddItem(emitButton, 0, 1, false)

	buttons = []*tview.Button{scanButton, backButton}
	pages.AddPage("infrared", IRFlex, true, false)
	enableTabFocus(IRFlex, buttons)
}

func loadIRScan() {
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

				// Stop LIRC services
				stopServices := []*exec.Cmd{
					exec.Command("sudo", "systemctl", "stop", "lircd"),
					exec.Command("sudo", "killall", "lircd"),
				}

				for _, cmd := range stopServices {
					if err := cmd.Run(); err != nil {
						app.QueueUpdateDraw(func() {
							scanText.SetText("[red]Error stopping services: " + err.Error())
						})
						return
					}
				}

				// Change infrared interface
				changeInterface := []*exec.Cmd{
					// Set driver to default
					exec.Command("sudo", "sed", "-i",
						"s/^driver.*/driver = default/",
						"/etc/lirc/lirc_options.conf"),

					// Set device to /dev/lirc0 (note escaped slashes)
					exec.Command("sudo", "sed", "-i",
						"s/^device.*/device = \\/dev\\/lirc0/",
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

				// Save to file
				content := strings.Join(timings, " ")
				if err := os.WriteFile("ir_timings.txt", []byte(content), 0644); err != nil {
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
	enableTabFocus(scanFlex, buttons)
}
