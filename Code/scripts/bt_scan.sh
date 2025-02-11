#!/bin/bash

# Create a timestamp for unique log file naming.
TIMESTAMP=$(date +%s)
LOG_FILE="/home/adamfitz395/Documents/GitHub/MultiTool-Project/Code/logfiles/btlogs/bluetooth_scan_$TIMESTAMP.log"

# Ensure the directory for logfile exists.
mkdir -p "$(dirname "$LOG_FILE")"
touch "$LOG_FILE"

# Create a temporary file to store the hcitool output.
TMP_LOG=$(mktemp)

DURATION=$1

# Restart Bluetooth to ensure no errors occur
echo "Starting Bluetooth scan..."
sudo hciconfig hci0 down
sudo hciconfig hci0 up

# Run hcitool lescan in the background
sudo hcitool lescan > "$TMP_LOG" 2>&1 &

SCAN_PID=$!  # Get the process ID of hcitool lescan

# Start a background process that waits x seconds, then sends signal to interrupt
( sleep "$DURATION"; sudo kill -SIGINT "$SCAN_PID" ) &

echo "Scanning for $DURATION seconds..."
wait "$SCAN_PID"  # Wait for hcitool lescan to stop

echo "Processing discovered devices..."

# Extract and log discovered devices
grep -E '([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}' "$TMP_LOG" | while read -r line; do
    mac=$(echo "$line" | awk '{print $1}')
    name=$(echo "$line" | cut -d ' ' -f2- | sed 's/^ *//')

    log_entry="$mac - $name" # 3A:1A:52:F2:65:4F - ET-4800 Series

    # Append the entry to the log file and print it.
    echo "$log_entry" | tee -a "$LOG_FILE"
done

# Clean up the temporary file.
rm "$TMP_LOG"

echo "Scan complete! Results saved to $LOG_FILE"
exit 0
