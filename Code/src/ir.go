package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"os/exec"
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
					scanText.SetText("[yellow]Stopping LIRC service...")
				})

				// Stop LIRC service
				if err := exec.Command("sudo", "systemctl", "stop", "lircd").Run(); err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Error stopping LIRC: " + err.Error())
					})
					return
				}

				// Capture IR signal
				app.QueueUpdateDraw(func() {
					scanText.SetText("[green]Press remote button now...")
				})

				cmd := exec.Command("mode2", "-d", "/dev/lirc1")
				output, err := cmd.CombinedOutput()
				if err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Capture failed: " + err.Error())
					})
					return
				}

				// Save raw data
				if err := os.WriteFile("ir_raw.txt", output, 0644); err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Save failed: " + err.Error())
					})
					return
				}

				app.QueueUpdateDraw(func() {
					scanText.SetText(fmt.Sprintf(
						"[green]Captured %d pulses\nSaved to [white]ir_raw.txt",
						len(output)))
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

func loadIREmit() {
	buttons = nil
	emitText := tview.NewTextView().
		SetDynamicColors(true).SetText("[green]Ready to scan!")
	emitText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	emitDropdown := tview.NewDropDown().
		AddOption("")
	emitDropdown.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	emitButton := tview.NewButton("Emit").
		SetSelectedFunc(func() {
			go func() {

			}()
		})
	emitButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	emitBackButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("infrared")
		})
	emitBackButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	emitFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(emitText, 0, 3, false).
		AddItem(emitDropdown, 0, 1, true).
		AddItem(emitButton, 0, 1, true).
		AddItem(emitBackButton, 0, 1, false)

	buttons = []*tview.Button{emitButton, emitBackButton}
	pages.AddPage("iremit", emitFlex, true, false)
	enableTabFocus(emitFlex, buttons)
}
