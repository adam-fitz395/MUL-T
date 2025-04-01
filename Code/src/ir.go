package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/smarthome-go/infrared"
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
				ifScanner := infrared.Scanner
				if err := ifScanner.Setup(4); err != nil {
					scanText.SetText(err.Error())
				}

				receivedCode, err := ifScanner.Scan()
				if err != nil {
					scanText.SetText(err.Error())
				}

				if receivedCode != "" {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[green]Found infrared code:[blue]" + receivedCode)
					})
				} else {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Unable to find infrared code.")
					})
				}
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
