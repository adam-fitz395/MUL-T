package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create a button with a visible border
	wifiButton := tview.NewButton("Wifi").
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetBackgroundColor(tcell.ColorWhite)

	// Create a Flex in row-direction (vertical stacking)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		// We'll give the button a fixed size of 3 rows (you can adjust),
		// proportion=1 just means "take leftover space if any is left",
		// and `true` means this item can receive focus.
		AddItem(wifiButton, 2, 1, true).
		SetTitle("Wifi").
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	// IMPORTANT: Use 'true' so the root fills the terminal window.
	if err := app.SetRoot(flex, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}
