#!/bin/bash

# Define variables
SERVICE_NAME="config-server"
BINARY_NAME="config-server"
INSTALL_DIR="/usr/local/bin"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"

# Step 1: Build the Go server
echo "Building the Go server..."
go build -o $BINARY_NAME ./server/server.go

# Step 2: Move the binary to the installation directory
echo "Installing the binary to $INSTALL_DIR..."
sudo mv $BINARY_NAME $INSTALL_DIR/

# Step 3: Create a systemd service file
echo "Creating systemd service file at $SERVICE_FILE..."
sudo bash -c "cat > $SERVICE_FILE" <<EOL
[Unit]
Description=Config Server
After=network.target

[Service]
ExecStart=$INSTALL_DIR/$BINARY_NAME
Restart=always
User=nobody
Group=nogroup
Environment=GO_ENV=production

[Install]
WantedBy=multi-user.target
EOL

# Step 4: Enable and start the service
echo "Enabling and starting the service..."
sudo systemctl enable $SERVICE_NAME
sudo systemctl start $SERVICE_NAME

echo "Installation complete. The service is now running."
