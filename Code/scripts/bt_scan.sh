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

# Run hcitool lescan with unbuffered output
TMP_LOG=$(mktemp)
sudo stdbuf -oL -eL hcitool lescan > "$TMP_LOG" 2>&1 &
SCAN_PID=$!

# Terminate after duration
( sleep "$DURATION"; sudo kill -SIGTERM "$SCAN_PID"; sleep 1 ) &

echo "Scanning for $DURATION seconds..."
wait "$SCAN_PID" 2>/dev/null

# Process results
echo "Processing discovered devices..."
grep -E '([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}' "$TMP_LOG" | while read -r line; do
  mac=$(echo "$line" | awk '{print $1}')
  name=$(echo "$line" | cut -d ' ' -f2- | sed 's/^ *//')
  echo "$mac - $name" | tee -a "$LOG_FILE"
done

rm "$TMP_LOG"
echo "Scan complete! Results saved to $LOG_FILE"
exit 0