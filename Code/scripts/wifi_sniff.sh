#!/bin/bash

TIMESTAMP=$(date +%s)
LOG_DIR="/home/adamfitz395/Documents/GitHub/MultiTool-Project/Code/logfiles/wifisnifflogs"
LOG_FILE="$LOG_DIR/sniff_log_$TIMESTAMP.pcap"
duration=$1

# Ensure directory exists and set permissions
sudo mkdir -p "$LOG_DIR"
sudo chown root:root "$LOG_DIR"
sudo chmod 755 "$LOG_DIR"

# Start scan using Pi interface and log PID for termination
sudo ettercap -T -w "$LOG_FILE" -i wlan0 &
SCAN_PID=$!

# Start background process to kill scan after set duration
(
  sleep "$duration"
  echo "Sending SIGINT to ettercap (PID: $SCAN_PID)"
  sudo kill -SIGINT "$SCAN_PID"
) &

echo "Scanning for $duration seconds..."
wait "$SCAN_PID" 2>/dev/null  # Suppress 'pid is not a child' errors
echo "Scan complete! Results saved to $LOG_FILE"
exit 0