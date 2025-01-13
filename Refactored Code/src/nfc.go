package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

// PN532 Command Definitions
const (
	CMD_GET_FIRMWARE_VERSION = 0x02
	PREAMBLE                 = 0x00
	STARTCODE1               = 0x00
	STARTCODE2               = 0xFF
	POSTAMBLE                = 0x00
)

func main() {
	// Open the UART port
	port, err := serial.OpenPort(&serial.Config{
		Name: "/dev/ttyS0", // Replace with your UART port (e.g., /dev/ttyUSB0 or COM3)
		Baud: 115200,
	})
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer port.Close()

	// Send GetFirmwareVersion command
	err = sendCommand(port, []byte{CMD_GET_FIRMWARE_VERSION})
	if err != nil {
		log.Fatalf("Failed to send command: %v", err)
	}

	// Wait for response
	response := make([]byte, 255)
	time.Sleep(100 * time.Millisecond) // Allow time for the PN532 to respond
	n, err := port.Read(response)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	if n > 0 {
		fmt.Printf("Received %d bytes: %X\n", n, response[:n])

		// Check if the response matches the expected format
		if response[6] == 0x02 { // Command echo for GetFirmwareVersion
			fmt.Println("PN532 is connected and responsive!")
		} else {
			fmt.Println("Unexpected response, PN532 may not be connected properly.")
		}
	} else {
		fmt.Println("No response received. Is the PN532 connected?")
	}
}

// sendCommand sends a command to the PN532
func sendCommand(port *serial.Port, data []byte) error {
	// Calculate checksum
	checksum := byte(0x00)
	for _, b := range data {
		checksum += b
	}
	checksum = ^checksum + 1

	// Build the full frame
	frame := []byte{PREAMBLE, STARTCODE1, STARTCODE2}
	length := byte(len(data))
	frame = append(frame, length, ^length+1)
	frame = append(frame, data...)
	frame = append(frame, checksum, POSTAMBLE)

	// Send the frame
	_, err := port.Write(frame)
	return err
}
