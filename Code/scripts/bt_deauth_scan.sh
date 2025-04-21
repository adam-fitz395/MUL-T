#!/bin/bash

TEMP_LOG_FILE="../logfiles/btlogs/tmpscan.log"

# Ensure directory exists
mkdir -p "$(dirname "$TEMP_LOG_FILE")"
touch "$LOG_FILE"

# Restart Bluetooth
echo "Starting Bluetooth scan..."
sudo hciconfig hci0 down
sudo hciconfig hci0 up

# Create a FIFO (named pipe) for real-time output
FIFO_LOG=$(mktemp -u)
mkfifo "$FIFO_LOG"

# Start hcitool lescan, writing to the FIFO
sudo timeout 5 stdbuf -oL hcitool lescan > "$FIFO_LOG" 2>&1 &

# Process FIFO output in the background
(
  declare -A seen_macs  # Track MAC addresses we've already seen
  while read -r line; do
    if echo "$line" | grep -qE '([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}'; then
      mac=$(echo "$line" | awk '{print $1}')
      name=$(echo "$line" | cut -d ' ' -f2- | sed 's/^ *//')

            # Skip devices with unknown name
            if [ "$name" == "(unknown)" ]; then
                continue
            fi

             # Only process if MAC hasn't been seen before
            if [[ -z "${seen_macs[$mac]+_}" ]]; then
              echo "$mac - $name" | tee -a "$LOG_FILE"
              seen_macs[$mac]=1  # Mark as seen AFTER processing
            fi
    fi
  done < "$FIFO_LOG"
) &

# Wait for the timeout to finish
wait

# Cleanup
rm -f "$FIFO_LOG"
echo "Scan complete! Results saved to $TEMP_LOG_FILE"
exit 0