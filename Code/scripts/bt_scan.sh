#!/bin/bash

# Timestamp for unique log file naming
TIMESTAMP=$(date +%s)
LOG_FILE="$HOME/Documents/GitHub/MultiTool-Project/Code/logfiles/btlogs/bluetooth_scan_$TIMESTAMP.log"

# Start scanning in the background using bluetoothctl
echo -e 'scan on' | bluetoothctl &
SCAN_PID=$!
sleep 5  # Adjust the scanning duration as needed

# Stop scanning
echo -e 'scan off' | bluetoothctl
wait $SCAN_PID  # Wait for the scan to finish

# Extract the device information from the 'bluetoothctl devices' output
echo -e 'devices' | bluetoothctl | while read -r line; do
    if [[ "$line" =~ ([0-9A-Fa-f:]{17}) ]]; then
        DEVICE_MAC=${BASH_REMATCH[1]}  # Extract the MAC address
        DEVICE_NAME=$(echo "$line" | cut -d ' ' -f 3-)  # Extract the device name (starting from the 3rd field)
        echo "$DEVICE_MAC - $DEVICE_NAME" | tee -a "$LOG_FILE"  # Append the MAC and Name to the log file
    fi
done

exit 0
