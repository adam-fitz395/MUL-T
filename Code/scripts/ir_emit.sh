#!/bin/bash

REMOTE_NAME="MY_REMOTE"
BUTTON_NAME="KEY_POWER"
CONFIG_DIR="ir_signals"
CONFIG_FILE="${CONFIG_DIR}/${REMOTE_NAME}.lircd.conf"

# Set LIRC options for the transmitter
sudo sed -i 's/^driver.*/driver = default/' /etc/lirc/lirc_options.conf
sudo sed -i 's/^device.*/device = \/dev\/lirc0/' /etc/lirc/lirc_options.conf

# Restart LIRC service
sudo systemctl restart lircd

# Ensure LIRC service is running
sudo systemctl start lircd

# Load the remote configuration
sudo cp ${CONFIG_FILE} /etc/lirc/lircd.conf
sudo systemctl restart lircd

# Send the IR signal
irsend SEND_ONCE ${REMOTE_NAME} ${BUTTON_NAME}

echo "IR signal sent using remote: ${REMOTE_NAME}, button: ${BUTTON_NAME}"
