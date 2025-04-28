#!/bin/bash

# File configuration
CONFIG_FILE="/etc/lirc/lirc_options.conf"
RECEIVER_DEVICE="/dev/lirc1"
OUTPUT_DIR="../logfiles/rawir"
MIN_TIMINGS=20
DURATION=$1

# Make  output directory if it does not exist
mkdir -p "$OUTPUT_DIR"

# Stop LIRC service with error handling
echo "Stopping LIRC services..."
sudo systemctl stop lircd 2>/dev/null
status=$?

# Handle systemctl exit codes
if [ $status -ne 0 ] && [ $status -ne 5 ]; then
  echo "Error: Failed to stop lircd service (code $status)" >&2
  exit 1
fi

# Update LIRC configuration to use receiver LIRC1
echo "Updating interface configuration..."
sudo sed -i "s/^driver.*/driver = default/" "$CONFIG_FILE" || exit 1
sudo sed -i "s/^device.*/device = \/dev\/lirc1/" "$CONFIG_FILE" || exit 1

# Capture IR signals
echo "Hold down remote button for $DURATION seconds..."
raw_output=$(timeout "$DURATION" mode2 -d "$RECEIVER_DEVICE" 2>&1)
status=$?

# Handle capture errors (ignore timeout)
if [ $status -ne 0 ] && [ $status -ne 124 ]; then
  echo "Capture error: $raw_output" >&2
  exit 1
fi

# Process timings
timings=()
while IFS= read -r line; do
  if [[ $line =~ (pulse|space)[[:space:]]+([0-9]+) ]]; then
    timings+=("${BASH_REMATCH[2]}")
  fi
done <<< "$raw_output"

# Validate results
if [ ${#timings[@]} -lt $MIN_TIMINGS ]; then
  echo "Error: Insufficient data (${#timings[@]} timings)" >&2
  exit 1
fi

# Fix odd timing count
if [ $(( ${#timings[@]} % 2 )) -ne 0 ]; then
  unset 'timings[${#timings[@]}-1]'
fi

# Create output directory
mkdir -p "$OUTPUT_DIR" || exit 1

# Generate filename
timestamp=$(date +%Y%m%d_%H%M%S)
output_file="$OUTPUT_DIR/ir_$timestamp.txt"

# Save to file
echo "${timings[*]}" > "$output_file" || exit 1

echo "Success: Captured ${#timings[@]} timings"
echo "Saved to: $output_file"
