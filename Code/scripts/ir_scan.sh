#!/bin/bash

LOG_DIR="../logfiles/irscanlogs"
REMOTE_NAME="MY_REMOTE"
BUTTON_NAME="KEY_POWER"
CONFIG_DIR="ir_signals"

# Set LIRC options for the receiver
sudo sed -i 's/^driver.*/driver = default/' /etc/lirc/lirc_options.conf
sudo sed -i 's/^device.*/device = \/dev\/lirc1/' /etc/lirc/lirc_options.conf

# Capture raw signal for 10 seconds
sudo systemctl stop lircd
sudo killall lircd 2>/dev/null
timeout 10 mode2 -d /dev/lirc1 > power_button.raw

# Process timings
grep -Eo '[0-9]+' power_button.raw | tr '\n' ' ' > power_button.timings

# Create LIRC config
GAP=$(head -n 1 power_button.timings | awk '{print $1}')
CODE=$(awk '{$1=""; print $0}' power_button.timings | sed 's/^ //')

# Write the configuration file with variable expansion
cat <<-EOF > ${CONFIG_DIR}/${REMOTE_NAME}.lircd.conf
begin remote
  name  MY_REMOTE
  flags RAW_CODES
  eps            30
  aeps          100
  gap          200000
  toggle_bit_mask 0x0

  begin raw_codes
    name ${BUTTON_NAME}
      ${CODE}
  end raw_codes
end remote
EOF

# Optional: Add frequency if using modulated signal
sed -i "/^begin remote/a \ \ frequency    38000" ${CONFIG_DIR}/${REMOTE_NAME}.lircd.conf

# Move files to organized locations
mkdir -p ${CONFIG_DIR}
mv power_button.* ${CONFIG_DIR}/
mv ${CONFIG_DIR}/${REMOTE_NAME}.lircd.conf ${CONFIG_DIR}/

echo "Config file created: ${CONFIG_DIR}/${REMOTE_NAME}.lircd.conf"
