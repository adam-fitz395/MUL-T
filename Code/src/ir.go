package main

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
	"os/exec"
)

// LoadIRMenu is a function that loads the menu for infrared interactions
func LoadIRMenu() {
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
	EnableTabFocus(IRFlex, buttons)
}

// LoadIRScan is a function that loads the infrared scanning sub-menu
func LoadIRScan() {
	buttons = nil
	var duration int

	scanText := tview.NewTextView().
		SetDynamicColors(true).SetText("[green]Ready to scan!")
	scanText.SetBorder(true).
		SetBorderColor(tcell.ColorWhite)

	scanDuration := tview.NewDropDown().SetLabel("Duration (Seconds): ").SetLabelColor(tcell.ColorWhite)
	scanDuration.
		AddOption("10",
			func() {
				duration = 10
			}).
		AddOption("20", func() {
			duration = 20
		}).
		AddOption("30", func() {
			duration = 30
		}).AddOption("60", func() {
		duration = 60
	})
	scanDuration.SetCurrentOption(0)

	scanButton := tview.NewButton("Scan").
		SetSelectedFunc(func() {
			go func() {
				// Start by clearing or setting initial text
				app.QueueUpdateDraw(func() {
					scanText.SetText("[yellow]Running capture script...\n")
				})

				durationStr := fmt.Sprintf("%d", duration)
				cmd := exec.Command("bash", "/path/to/your/script.sh", durationStr) // Update path!

				// Get stdout pipe
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Failed to get stdout: " + err.Error())
					})
					return
				}

				// Get stderr pipe (optional, if you want error messages too)
				stderr, err := cmd.StderrPipe()
				if err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Failed to get stderr: " + err.Error())
					})
					return
				}

				if err := cmd.Start(); err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Failed to start script: " + err.Error())
					})
					return
				}

				// Merge stdout and stderr into one scanner
				outputReader := io.MultiReader(stdout, stderr)
				scanner := bufio.NewScanner(outputReader)

				for scanner.Scan() {
					line := scanner.Text()

					// Update scanText with each new line
					app.QueueUpdateDraw(func() {
						currentText := scanText.GetText(true) // 'true' to get without markup parsing if needed
						scanText.SetText(currentText + "\n" + line)
					})
				}

				if err := scanner.Err(); err != nil {
					app.QueueUpdateDraw(func() {
						scanText.SetText("[red]Error reading output: " + err.Error())
					})
					return
				}

				// Wait for script to finish
				if err := cmd.Wait(); err != nil {
					app.QueueUpdateDraw(func() {
						currentText := scanText.GetText(true)
						scanText.SetText(currentText + "\n[red]Script finished with error: " + err.Error())
					})
					return
				}

				// Done successfully
				app.QueueUpdateDraw(func() {
					currentText := scanText.GetText(true)
					scanText.SetText(currentText + "\n[green]Script completed successfully.")
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
		AddItem(scanDuration, 0, 1, false).
		AddItem(scanButton, 0, 1, true).
		AddItem(scanBackButton, 0, 1, false)

	buttons = []*tview.Button{scanButton, scanBackButton}
	pages.AddPage("irscan", scanFlex, true, false)
	EnableTabFocus(scanFlex, buttons)
}
