#!/bin/bash

# Define variables
SERVICE_NAME="config-server"
BINARY_NAME="config-server"
INSTALL_DIR="/usr/local/bin"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"

# Step 1: Stop the service
echo "Stopping the service..."
sudo systemctl stop $SERVICE_NAME

# Step 2: Disable the service
echo "Disabling the service..."
sudo systemctl disable $SERVICE_NAME

# Step 3: Remove the systemd service file
echo "Removing systemd service file at $SERVICE_FILE..."
sudo rm -f $SERVICE_FILE

# Step 4: Remove the binary
echo "Removing the binary from $INSTALL_DIR..."
sudo rm -f $INSTALL_DIR/$BINARY_NAME

echo "Uninstallation complete. The service has been removed."
