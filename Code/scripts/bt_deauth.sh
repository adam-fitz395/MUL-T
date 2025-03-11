#!/bin/bash

LOG_DIR="/home/adamfitz395/Documents/GitHub/MultiTool-Project/Code/logfiles/btlogs"
PID_FILE="$LOG_DIR/deauth.pid"

sudo hciconfig hci0 down
sudo hciconfig hci0 up

sudo timeout 10 hcitool lescan