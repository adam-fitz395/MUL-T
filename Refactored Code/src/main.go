package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app     *tview.Application
	pages   *tview.Pages
	buttons []*tview.Button
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
	buttons = nil // Set buttons to nil to clear buttons from previous menu

	wifiButton := tview.NewButton("WiFi Menu").
		SetSelectedFunc(func() {
			pages.SwitchToPage("wifi") // Switch to the Wi-Fi page
		})

	buttons = append(buttons, wifiButton)

	exitButton := tview.NewButton("Exit").
		SetSelectedFunc(func() {
			app.Stop() // Exit the application
		})

	buttons = append(buttons, exitButton)

	mainFlex := tview.NewFlex().
		AddItem(wifiButton, 0, 1, true).
		AddItem(exitButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	pages.AddPage("main", mainFlex, true, true) // Add the main page to pages
	enableTabFocus(mainFlex, buttons)
}

// Function that allows user to switch focus between elements using the "Tab" button
func enableTabFocus(layout *tview.Flex, focusables []*tview.Button) {
	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			// Find the currently focused button and move to the next one
			for i, btn := range focusables {
				if app.GetFocus() == btn {
					// Cycle to the next button
					nextIndex := (i + 1) % len(focusables)
					app.SetFocus(focusables[nextIndex])
					break
				}
			}
			return nil // Prevent default Tab behavior
		}
		return event
	})
}
