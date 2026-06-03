#!/bin/bash

# GoSync Server - Install as System Service (systemd)
# This script installs GoSync as a system service on Linux

set -e

echo "=== GoSync Server - Service Installation ==="
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (use sudo)" 
   exit 1
fi

# Get the current user who invoked sudo
ACTUAL_USER=${SUDO_USER:-$USER}
CURRENT_DIR=$(pwd)

read -p "Install directory (default: /opt/gosync): " INSTALL_DIR
INSTALL_DIR=${INSTALL_DIR:-/opt/gosync}

echo "Creating installation directory..."
mkdir -p "$INSTALL_DIR"

echo "Copying server files..."
cp -r server/* "$INSTALL_DIR/"
cd "$INSTALL_DIR"

# Build if binary doesn't exist
if [[ ! -f "gosync-server" ]]; then
    echo "Building server binary..."
    go build -o gosync-server main.go config.go
fi

# Create systemd service file
echo "Creating systemd service..."
cat > /etc/systemd/system/gosync.service << SERVICEEOF
[Unit]
Description=GoSync Server
After=network.target

[Service]
Type=simple
User=$ACTUAL_USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/gosync-server
Restart=always
RestartSec=10

# Security settings
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
SERVICEEOF

echo "Reloading systemd daemon..."
systemctl daemon-reload

echo ""
echo "=== Installation Complete! ==="
echo ""
echo "Service commands:"
echo "  Start:   sudo systemctl start gosync"
echo "  Stop:    sudo systemctl stop gosync"
echo "  Restart: sudo systemctl restart gosync"
echo "  Status:  sudo systemctl status gosync"
echo "  Enable auto-start: sudo systemctl enable gosync"
echo "  View logs: sudo journalctl -u gosync -f"
echo ""

read -p "Enable auto-start on boot? (y/n): " AUTO_START
if [[ "$AUTO_START" == "y" || "$AUTO_START" == "Y" ]]; then
    systemctl enable gosync
    echo "Auto-start enabled!"
fi

read -p "Start service now? (y/n): " START_NOW
if [[ "$START_NOW" == "y" || "$START_NOW" == "Y" ]]; then
    systemctl start gosync
    echo "Service started!"
    echo ""
    systemctl status gosync
fi
