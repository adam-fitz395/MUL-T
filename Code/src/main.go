// Package main provides a terminal-based user interface (TUI) for wireless pen-testing tools.
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

// main initializes and runs the TUI application.
// It sets up the application instance, creates page container,
// initializes all menus, and starts the event loop.
func main() {
	app = tview.NewApplication()
	pages = tview.NewPages()

	InstantiateMenus()

	// Run application or return error
	if err := app.SetRoot(pages, true).
		SetFocus(pages).
		EnableMouse(true).
		Run(); err != nil {
		fmt.Print("Error running app:", err)
	}
}

// InstantiateMenus is a function initializes all application menus and pages.
func InstantiateMenus() {
	LoadMainMenu()
	LoadWifiMenu()
	LoadSniffingMenu()
	LoadScanMenu()
	LoadMITMMenu()
	LoadBluetoothMenu()
	LoadBluetoothScan()
	LoadBluetoothDeauth()
	LoadIRMenu()
	LoadIRScan()
}

// LoadMainMenu is a function that creates and configures the root menu page.
func LoadMainMenu() {
	buttons = nil // Set buttons to nil to clear buttons from previous menu

	wifiButton := tview.NewButton("Wi-Fi").
		SetSelectedFunc(func() {
			pages.SwitchToPage("wifi") // Switch to the Wi-Fi page
		})
	wifiButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)
	wifiButton.SetBackgroundColorActivated(tcell.ColorGreen).
		SetLabelColorActivated(tcell.ColorWhite)

	bluetoothButton := tview.NewButton("Bluetooth").
		SetSelectedFunc(func() {
			pages.SwitchToPage("bluetooth")
		})
	bluetoothButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)
	bluetoothButton.SetBackgroundColorActivated(tcell.ColorBlue).
		SetLabelColorActivated(tcell.ColorWhite)

	IRButton := tview.NewButton("Infrared").
		SetSelectedFunc(func() {
			pages.SwitchToPage("infrared")
		})
	IRButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)
	IRButton.SetBackgroundColorActivated(tcell.ColorRed).
		SetLabelColorActivated(tcell.ColorWhite)

	exitButton := tview.NewButton("Exit").
		SetSelectedFunc(func() {
			app.Stop() // Exit the application
		})
	exitButton.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	mainFlex := tview.NewFlex().
		AddItem(wifiButton, 0, 1, true).
		AddItem(bluetoothButton, 0, 1, true).
		AddItem(IRButton, 0, 1, false).
		AddItem(exitButton, 0, 1, false).
		SetDirection(tview.FlexRow)

	buttons = []*tview.Button{wifiButton, bluetoothButton, IRButton, exitButton}
	pages.AddPage("main", mainFlex, true, true) // Add the main page to pages
	EnableTabFocus(mainFlex, buttons)
}

// EnableTabFocus is a function that allows users to switch focus between elements using the "Tab" button
// Parameters:
//   - layout: The Flex container to enable navigation in
//   - focusables: Slice of buttons that should receive focus
func EnableTabFocus(layout *tview.Flex, focusables []*tview.Button) {
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
