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
		fmt.Print("Stinky code bad")
	}
}

// Function that instantiates all menus
func instantiateMenus() {
	loadMainMenu()
	loadWifiMenu()
	loadSniffingMenu()
	loadScanMenu()
	loadIRMenu()
}

// Function to load main menu
func loadMainMenu() {
	buttons = nil // Set buttons to nil to clear buttons from previous menu

	wifiButton := tview.NewButton("Wi-Fi Menu").
		SetSelectedFunc(func() {
			pages.SwitchToPage("wifi") // Switch to the Wi-Fi page
		})

	wifiButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	wifiButton.SetBackgroundColorActivated(tcell.ColorGreen).
		SetLabelColorActivated(tcell.ColorWhite)

	IRButton := tview.NewButton("Infrared").
		SetSelectedFunc(func() {
			pages.SwitchToPage("infrared")
		})

	IRButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)
	IRButton.SetBackgroundColorActivated(tcell.ColorRed).
		SetLabelColorActivated(tcell.ColorWhite)

	NFCButton := tview.NewButton("NFC").
		SetSelectedFunc(func() {
			pages.SwitchToPage("nfc")
			app.Stop()
			UIDScan()
		})

	NFCButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)
	NFCButton.SetBackgroundColorActivated(tcell.ColorBlueViolet).
		SetLabelColorActivated(tcell.ColorWhite)

	exitButton := tview.NewButton("Exit").
		SetSelectedFunc(func() {
			app.Stop() // Exit the application
		})

	exitButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	mainFlex := tview.NewFlex().
		AddItem(wifiButton, 0, 1, true).
		AddItem(IRButton, 0, 1, false).
		AddItem(NFCButton, 0, 1, false).
		AddItem(exitButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = []*tview.Button{wifiButton, IRButton, NFCButton, exitButton}
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
