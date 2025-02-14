#!/bin/bash

TIMESTAMP=$(date +%s)
LOG_FILE="/home/adamfitz395/Documents/GitHub/MultiTool-Project/Code/logfiles/btlogs/bluetooth_scan_$TIMESTAMP.log"
DURATION=$1

# Ensure directory exists
mkdir -p "$(dirname "$LOG_FILE")"
touch "$LOG_FILE"

# Restart Bluetooth
echo "Starting Bluetooth scan..."
sudo hciconfig hci0 down
sudo hciconfig hci0 up

# Create a FIFO (named pipe) for real-time output
FIFO_LOG=$(mktemp -u)
mkfifo "$FIFO_LOG"

# Start hcitool lescan, writing to the FIFO
sudo timeout "$DURATION" stdbuf -oL hcitool lescan > "$FIFO_LOG" 2>&1 &

# Process FIFO output in the background
(
  while read -r line; do
    if echo "$line" | grep -qE '([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}'; then
      mac=$(echo "$line" | awk '{print $1}')
      name=$(echo "$line" | cut -d ' ' -f2- | sed 's/^ *//')
      echo "$mac - $name" | tee -a "$LOG_FILE"
    fi
  done < "$FIFO_LOG"
) &

# Wait for the timeout to finish
wait

# Cleanup
rm -f "$FIFO_LOG"
echo "Scan complete! Results saved to $LOG_FILE"
exit 0