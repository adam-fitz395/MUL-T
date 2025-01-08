package main

import (
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	pages *tview.Pages
)

func main() {
	app = tview.NewApplication()
	pages = tview.NewPages()

	loadMainMenu()
	loadWifiMenu()

	// IMPORTANT: Use 'true' so the root fills the terminal window.
	if err := app.SetRoot(pages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func loadMainMenu() {
	wifiButton := tview.NewButton("WiFi Menu").
		SetSelectedFunc(func() {
			pages.SwitchToPage("wifi") // Switch to the WiFi page
		})

	exitButton := tview.NewButton("Exit").
		SetSelectedFunc(func() {
			app.Stop() // Exit the application
		})

	mainFlex := tview.NewFlex().
		AddItem(wifiButton, 0, 1, true).
		AddItem(exitButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	pages.AddPage("main", mainFlex, true, true) // Add the main page to pages
}
