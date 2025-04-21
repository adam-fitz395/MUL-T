#!/bin/bash

LOG_DIR="../logfiles/irscanlogs"
REMOTE_NAME="MY_REMOTE"
BUTTON_NAME="KEY_POWER"

# Capture raw signal for 10 seconds
sudo systemctl stop lircd
sudo killall lircd 2>/dev/null
timeout 10 mode2 -d /dev/lirc1 > power_button.raw

# Process timings
grep -Eo '[0-9]+' power_button.raw | tr '\n' ' ' > power_button.timings

# Create LIRC config
GAP=$(head -n 1 power_button.timings | awk '{print $1}')
CODE=$(awk '{$1=""; print $0}' power_button.timings | sed 's/^ //')

cat << EOF > ${REMOTE_NAME}.lircd.conf
begin remote
  name  $REMOTE_NAME
  flags RAW_CODES
  eps            30
  aeps          100
  gap          200000  # Default gap if none detected
  toggle_bit_mask 0x0

  begin raw_codes
    name $BUTTON_NAME
      $CODE
  end raw_codes
end remote
EOF

# Optional: Add frequency if using modulated signal
sed -i "/^begin remote/a \ \ frequency    38000" ${REMOTE_NAME}.lircd.conf

# Move files to organized locations
mkdir -p ir_signals
mv power_button.* ir_signals/
mv ${REMOTE_NAME}.lircd.conf ir_signals/

echo "Config file created: ir_signals/${REMOTE_NAME}.lircd.conf"
