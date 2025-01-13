package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
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
	command := []byte{0x00, 0x00, 0xFF, 0x02, 0xFE, 0xD4, 0x02, 0x2A, 0x00}
	_, err = port.Write(command)
	if err != nil {
		log.Fatalf("Failed to send command: %v", err)
	}

	// Wait for response
	time.Sleep(100 * time.Millisecond) // Wait for PN532 to respond
	response := make([]byte, 255)
	n, err := port.Read(response)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Check if we received a response
	if n > 0 {
		fmt.Printf("PN532 is connected! Received response: %X\n", response[:n])
	} else {
		fmt.Println("No response received. Is the PN532 connected?")
	}
}
