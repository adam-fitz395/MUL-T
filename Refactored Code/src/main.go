package main

import (
	"fmt"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	pages *tview.Pages
)

func main() {
	app = tview.NewApplication()
	pages = tview.NewPages()

	instantiateMenus()

	// Run application or return error
	if err := app.SetRoot(pages, true).
		SetFocus(pages).
		EnableMouse(true).
		Run(); err != nil {
		fmt.Println("An error has occured:", err)
		panic(err)
	}
}

// Function that instantiates all menus
func instantiateMenus() {
	loadMainMenu()
	loadWifiMenu()
	loadSniffingMenu()
}

// Function to load main menu
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
