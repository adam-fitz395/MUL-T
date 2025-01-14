package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func UIDScan() {
	// Open the UART port
	port, err := serial.OpenPort(&serial.Config{
		Name: "/dev/ttyS0", // Replace with your UART port (e.g., /dev/ttyUSB0 or COM3)
		Baud: 115200,
	})
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	fmt.Println("Opened serial port")
	defer port.Close()

	// Send GetFirmwareVersion command
	command := []byte{0x00, 0x00, 0xFF, 0x02, 0xFE, 0xD4, 0x02, 0x2A, 0x00}
	_, err = port.Write(command)
	if err != nil {
		log.Fatalf("Failed to send command: %v", err)
	}
	fmt.Println("Sent command")

	// Wait for response
	time.Sleep(1 * time.Second) // Wait for PN532 to respond
	response := make([]byte, 255)
	fmt.Println(response)
	n, err := port.Read(response)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	fmt.Printf("Response: %v\n", string(response[:n]))

	// Check if we received a response
	if n > 0 {
		fmt.Printf("PN532 is connected! Received response: %X\n", response[:n])
	} else {
		fmt.Println("No response received. Is the PN532 connected?")
	}
}
