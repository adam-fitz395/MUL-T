#!/bin/bash


TIMESTAMP=$(date +%s)
LOG_DIR="/home/adamfitz395/Documents/GitHub/MultiTool-Project/Code/logfiles/wifisnifflogs"
LOG_FILE="$LOG_DIR/sniff_log_$TIMESTAMP.pcap"
duration=$1

# Ensure directory exists and is root-owned
sudo mkdir -p "$LOG_DIR"
sudo chown root:root "$LOG_DIR"  # <-- Critical line
sudo chmod 755 "$LOG_DIR"

# Verify interface exists
if ! ip link show wlo1 &>/dev/null; then
  echo "Error: Wireless interface not found!"
  exit 1
fi

# Start ettercap
sudo ettercap -T -w "$LOG_FILE" -i wlo1 &

SCAN_PID=$!

(
  sleep "$duration"
  sudo kill -SIGINT "$SCAN_PID"
  ) &

echo "Scanning for $duration seconds..."
wait "$SCAN_PID"

FINAL_DIR="/home/adamfitz395/Documents/.../wifisnifflogs"
mv "$LOG_FILE" "$FINAL_DIR/"

echo "Scan complete! Results saved to $FINAL_DIR"
exit 0