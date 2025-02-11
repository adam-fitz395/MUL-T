#!/bin/bash

TIMESTAMP=$(date +%s)
LOG_FILE="/home/adamfitz395/Documents/GitHub/MultiTool-Project/Code/logfiles/wifisnifflogs/sniff_log_$TIMESTAMP.pcap"

duration=$1

sudo ettercap -T -w "$LOG_FILE" -i wlo1 &

SCAN_PID=$!

( sleep "$duration"; sudo kill -SIGINT "$SCAN_PID" ) & # Create background process and wait to send kill signal

echo "Scanning for 10 seconds..."
wait "$SCAN_PID"  # Wait for ettercap scan to stop

echo "Scan complete! Results saved to $LOG_FILE"
exit 0

