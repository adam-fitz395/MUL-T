package main

import (
	"bufio"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os/exec"
)

func main() {
	multi_menu := app.New()
	multi_window := multi_menu.NewWindow("Network Sniffer")

	// Label to display output
	outputLabel := widget.NewLabel("Output will appear here")

	multi_window.SetContent(container.NewVBox(
		outputLabel,
		widget.NewButton("Start Sniffing", func() {
			outputLabel.SetText("Starting sniffing...")

			// Run the command
			cmd := exec.Command("sudo", "ettercap", "-T", "-i", "wlo1")

			// Get a pipe to stdout
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Println("Error creating stdout pipe:", err)
				outputLabel.SetText("Error: " + err.Error())
				return
			}

			// Start the command
			if err := cmd.Start(); err != nil {
				log.Println("Error starting command:", err)
				outputLabel.SetText("Error: " + err.Error())
				return
			}

			// Read stdout line by line and update the label
			scanner := bufio.NewScanner(stdout)
			go func() {
				for scanner.Scan() {
					line := scanner.Text()
					log.Println("Output:", line) // Log the output for debugging
					outputLabel.SetText(line)    // Update the label with each line
				}
				if err := scanner.Err(); err != nil {
					log.Println("Error reading output:", err)
				}
			}()

			// Wait for the command to finish
			if err := cmd.Wait(); err != nil {
				log.Println("Command finished with error:", err)
				outputLabel.SetText("Command finished with error: " + err.Error())
			}
		}),
	))

	multi_window.ShowAndRun()
}
