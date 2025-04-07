package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
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
				// Initialize host driver
				app.QueueUpdateDraw(func() {
					scanText.SetText("[yellow]Initializing hardware...")
				})
				if _, err := host.Init(); err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Error initializing: " + err.Error())
					})
					return
				}

				// GPIO setup
				app.QueueUpdateDraw(func() {
					scanText.SetText("[yellow]Configuring GPIO...")
				})
				pin := gpioreg.ByName("GPIO17")
				if pin == nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]GPIO17 not found")
					})
					return
				}

				if err := pin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]GPIO error: " + err.Error())
					})
					return
				}

				// Signal capture
				app.QueueUpdateDraw(func() {
					scanText.SetText("[green]Press IR remote button now...")
				})
				if !pin.WaitForEdge(10 * time.Second) {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]No signal detected")
					})
					return
				}

				var durations []time.Duration
				var states []gpio.Level
				lastEdge := time.Now()
				const timeout = 100 * time.Millisecond

				for {
					if !pin.WaitForEdge(timeout) {
						break
					}
					now := time.Now()
					newState := pin.Read()
					durations = append(durations, now.Sub(lastEdge))
					states = append(states, !newState)
					lastEdge = now
				}

				if len(durations) == 0 {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]No signal captured")
					})
					return
				}

				// Save to file
				filename := "ir_signal.txt"
				file, err := os.Create(filename)
				if err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]File error: " + err.Error())
					})
					return
				}
				defer file.Close()

				for i, d := range durations {
					_, err := fmt.Fprintf(file, "%d %d\n", d.Microseconds(), states[i])
					if err != nil {
						app.QueueUpdateDraw(func() {
							scanText.SetText("[red]Write error: " + err.Error())
						})
						return
					}
				}

				app.QueueUpdateDraw(func() {
					scanText.SetText(fmt.Sprintf(
						"[green]Captured %d pulses\nSaved to [white]%s",
						len(durations), filename))
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
