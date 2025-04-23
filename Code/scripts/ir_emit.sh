#!/bin/bash

REMOTE_NAME="MY_REMOTE"
BUTTON_NAME="KEY_POWER"
CONFIG_DIR="ir_signals"
CONFIG_FILE="${CONFIG_DIR}/${REMOTE_NAME}.lircd.conf"

# Set LIRC options for the transmitter
sudo sed -i 's/^driver.*/driver = default/' /etc/lirc/lirc_options.conf
sudo sed -i 's/^device.*/device = \/dev\/lirc0/' /etc/lirc/lirc_options.conf

# Copy the configuration file to the correct location
sudo cp ${CONFIG_FILE} /etc/lirc/lircd.conf

# Restart LIRC service
sudo systemctl restart lircd

# Verify the remote name is recognized
if ! irsend LIST "" "" | grep -q "${REMOTE_NAME}"; then
  echo "Error: Remote '${REMOTE_NAME}' not recognized. Check the configuration file."
  exit 1
fi

# Send the IR signal
irsend SEND_ONCE ${REMOTE_NAME} ${BUTTON_NAME}

if [ $? -eq 0 ]; then
  echo "IR signal sent using remote: ${REMOTE_NAME}, button: ${BUTTON_NAME}"
else
  echo "Error: Failed to send IR signal. Check the IR transmitter hardware."
fi
