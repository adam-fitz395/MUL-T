package main

import (
	"github.com/rivo/tview"
)

// Function that loads the menu for infrared functions
func loadIRMenu() {
	buttons = nil

	scanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			pages.SwitchToPage("irscan")
		})

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("main")
		})

	IRFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(scanButton, 0, 1, true).
		AddItem(backButton, 0, 1, false)

	buttons = []*tview.Button{scanButton, backButton}
	pages.AddPage("infrared", IRFlex, true, false)
	enableTabFocus(IRFlex, buttons)
}
