package main

import "github.com/rivo/tview"

func loadWifiMenu() {
	sniffButton := tview.NewButton("Sniff").
		SetSelectedFunc(func() {
			app.Stop() // Logic for Sniff button
		})

	backButton := tview.NewButton("Back").
		SetSelectedFunc(func() {
			pages.SwitchToPage("main") // Switch back to the main page
		})

	wifiFlex := tview.NewFlex().
		AddItem(sniffButton, 0, 1, true).
		AddItem(backButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	pages.AddPage("wifi", wifiFlex, true, false) // Add the WiFi page to pages
}
